#include "llvm/ADT/APFloat.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/IR/BasicBlock.h"
#include "llvm/IR/Constants.h"
#include "llvm/IR/DerivedTypes.h"
#include "llvm/IR/Function.h"
#include "llvm/IR/IRBuilder.h"
#include "llvm/IR/LLVMContext.h"
#include "llvm/IR/Module.h"
#include "llvm/IR/Type.h"
#include "llvm/IR/Verifier.h"
#include <iostream>
#include <memory>
#include <vector>
#include <map>

enum Token {
	tok_eof = -1,

	// comands
	tok_def = -2,
	tok_extern = -3,

	// primary
	tok_identifier = -4,
	tok_number = -5,
};

static std::string IdentifierStr;
static double NumVal;

static int gettok() {
	static int LastChar = ' ';

	while (isspace(LastChar))
		LastChar = getchar();

	if (isalpha(LastChar)) {
		IdentifierStr = LastChar;
		while (isalnum((LastChar = getchar())))
			IdentifierStr += LastChar;

		if (IdentifierStr == "def")
			return tok_def;
		if (IdentifierStr == "extern")
			return tok_extern;
		return tok_identifier;
	}

	if (isdigit(LastChar) || LastChar == '.') {
		std::string NumStr;
		do {
			NumStr += LastChar;
			LastChar = getchar();
		} while (isdigit(LastChar) || LastChar == '.');

		NumVal = strtod(NumStr.c_str(), 0);
		return tok_number;
	}

	if (LastChar == '#') {
		// Comment until end of line
		do
			LastChar = getchar();
		while (LastChar != EOF && LastChar != '\n' && LastChar != '\r');

		if (LastChar != EOF)
			return gettok();
	}

	if (LastChar == EOF)
		return tok_eof;

	int ThisChar = LastChar;
	LastChar = getchar();
	return ThisChar;
}

class ExprAST {
public:
	virtual ~ExprAST() = default;
	virtual llvm::Value *codegen() = 0;
};

class NumberExprAST : public ExprAST {
	double Val;

public:
	NumberExprAST(double Val) : Val(Val) {}
  llvm::Value *codegen() override;
};

class VariableExprAST : public ExprAST {
	std::string Name;

public:
	VariableExprAST(const std::string &Name) : Name(Name) {}
  llvm::Value *codegen() override;
};

class BinaryExprAST : public ExprAST {
	char Op;
  std::unique_ptr<ExprAST> LHS, RHS;

public:
	BinaryExprAST(char Op, std::unique_ptr<ExprAST> LHS, std::unique_ptr<ExprAST> RHS)
		: Op(Op), LHS(std::move(LHS)), RHS(std::move(RHS)) {}
	llvm::Value *codegen() override;
};

class CallExprAST : public ExprAST {
	std::string Callee;
	std::vector<std::unique_ptr<ExprAST>> Args;

public:
	CallExprAST(const std::string &Callee, std::vector<std::unique_ptr<ExprAST>> Args)
		: Callee(Callee), Args(std::move(Args)) {}
	llvm::Value *codegen() override;
};

class PrototypeAST {
	std::string Name;
	std::vector<std::string> Args;

public:
	PrototypeAST(const std::string &Name, std::vector<std::string> Args)
		: Name(Name), Args(std::move(Args)) {}
	llvm::Function *codegen();
	std::string getName();
};

class FunctionAST {
	std::unique_ptr<PrototypeAST> Proto;
	std::unique_ptr<ExprAST> Body;

public:
	FunctionAST(std::unique_ptr<PrototypeAST> Proto, std::unique_ptr<ExprAST> Body)
		: Proto(std::move(Proto)), Body(std::move(Body)) {}
	llvm::Function *codegen();
};

static int CurTok;
static int getNextToken() {
	return CurTok = gettok();
}

std::unique_ptr<ExprAST> LogError(const char *Str) {
	fprintf(stderr, "Error: %s\n", Str);
	return nullptr;
}

std::unique_ptr<PrototypeAST> LogErrorP(const char *Str) {
	LogError(Str);
	return nullptr;
}

static std::unique_ptr<ExprAST> ParseExpression();

static std::unique_ptr<ExprAST> ParseNumberExpr() {
	auto Result = std::make_unique<NumberExprAST>(NumVal);
	getNextToken(); // consume the number
	return std::move(Result);
}

static std::unique_ptr<ExprAST> ParseParenExpr() {
	getNextToken(); // eat (
	auto V = ParseExpression();
	if (!V)
		return nullptr;

	if (CurTok != ')')
		return LogError("expected ')'");

	getNextToken(); // eat )
	return V;
}

static std::unique_ptr<ExprAST> ParseIdentifierExpr() {
	std::string IdName = IdentifierStr;

	getNextToken(); // eat identifier

	if (CurTok != '(') // Simple variable ref
		return std::make_unique<VariableExprAST>(IdName);

	getNextToken(); // eat (
	std::vector<std::unique_ptr<ExprAST>> Args;
	if (CurTok != ')') {
		while (true) {
			if (auto Arg = ParseExpression())
				Args.push_back(std::move(Arg));
			else
				return nullptr;

			if (CurTok == ')')
				break;

			if (CurTok != ',')
				return LogError("Expected ')' or ',' in argument list");

			getNextToken();
		}
	}

	getNextToken(); // Eat the ')'
	return std::make_unique<CallExprAST>(IdName, std::move(Args));
}

static std::unique_ptr<ExprAST> ParsePrimary() {
	switch (CurTok) {
	default:
		return LogError("unknown token when expecting an expression");
	case tok_identifier:
		return ParseIdentifierExpr();
	case tok_number:
		return ParseNumberExpr();
	case '(':
		return ParseParenExpr();
	}
}

static std::map<char, int> BinopPrecedence;

static int GetTokPrecedence() {
	if (!isascii(CurTok))
		return -1;

	int TokPrec = BinopPrecedence[CurTok];
	if (TokPrec <= 0) return -1;
	return TokPrec;
}

static std::unique_ptr<ExprAST> ParseBinOpRHS(int ExprPrec, std::unique_ptr<ExprAST> LHS) {
	while (true) {
		int TokPrec = GetTokPrecedence();

		if (TokPrec < ExprPrec)
			return LHS;

		int BinOp = CurTok;
		getNextToken();

		auto RHS = ParsePrimary();
		if (!RHS)
			return nullptr;

		int NextPrec = GetTokPrecedence();
		if (TokPrec < NextPrec) {
			RHS = ParseBinOpRHS(TokPrec + 1, std::move(RHS));
			if (!RHS)
				return nullptr;
		}

		LHS = std::make_unique<BinaryExprAST>(BinOp, std::move(LHS), std::move(RHS));
	}
}

static std::unique_ptr<ExprAST> ParseExpression() {
	auto LHS = ParsePrimary();
	if (!LHS)
		return nullptr;

	return ParseBinOpRHS(0, std::move(LHS));
}

static std::unique_ptr<PrototypeAST> ParsePrototype() {
	if (CurTok != tok_identifier)
		return LogErrorP("Expected function name in prototype");

	std::string FnName = IdentifierStr;
	getNextToken();

	if (CurTok != '(')
		return LogErrorP("Expected '(' in prototype");

	std::vector<std::string> ArgNames;
	while (getNextToken() == tok_identifier)
		ArgNames.push_back(IdentifierStr);
	if (CurTok != ')')
		return LogErrorP("Expected ')' in prototype");

	getNextToken(); // eat ')'
	return std::make_unique<PrototypeAST>(FnName, std::move(ArgNames));
}

static std::unique_ptr<FunctionAST> ParseDefinition() {
	getNextToken(); // eat def
	auto Proto = ParsePrototype();
	if (!Proto) return nullptr;

	if (auto E = ParseExpression())
		return std::make_unique<FunctionAST>(std::move(Proto), std::move(E));

	return nullptr;
}

static std::unique_ptr<PrototypeAST> ParseExtern() {
	getNextToken(); // eat extern
	return ParsePrototype();
}

static std::unique_ptr<FunctionAST> ParseTopLevelExpr() {
	if (auto E = ParseExpression()) {
		auto Proto = std::make_unique<PrototypeAST>("", std::vector<std::string>());
		return std::make_unique<FunctionAST>(std::move(Proto), std::move(E));
	}

	return nullptr;
}

static void HandleDefinition() {
	if (ParseDefinition()) {
		fprintf(stderr, "Parsed a function definition.\n");
	} else {
		getNextToken();
	}
}

static void HandleExtern() {
	if (ParseExtern()) {
		fprintf(stderr, "Parsed an extern\n");
	} else {
		getNextToken();
	}
}

static void HandleTopLevelExpression() {
	if (ParseTopLevelExpr()) {
		fprintf(stderr, "Parsed a top-level expr\n");
	} else {
		getNextToken();
	}
}

static void MainLoop() {
	while (true) {
		fprintf(stderr, "ready> ");
		switch (CurTok) {
		case tok_eof:
			return;
		case ';': // ignore top-level semicolons
			getNextToken();
			break;
		case tok_def:
			HandleDefinition();
			break;
		case tok_extern:
			HandleExtern();
			break;
		default:
			HandleTopLevelExpression();
			break;
		}
	}
}

static std::unique_ptr<llvm::LLVMContext> TheContext;
static std::unique_ptr<llvm::IRBuilder<>> Builder;
static std::unique_ptr<llvm::Module> TheModule;
static std::map<std::string, llvm::Value *> NamedValues;

llvm::Value *LogErrorV(const char *Str) {
	LogError(Str);
	return nullptr;
}

llvm::Value *NumberExprAST::codegen() {
	return llvm::ConstantFP::get(*TheContext, llvm::APFloat(Val));
}

llvm::Value *VariableExprAST::codegen() {
	llvm::Value *V = NamedValues[Name];
	if (!V)
		LogErrorV("Unknown variable name");

	return V;
}

llvm::Value *BinaryExprAST::codegen() {
	llvm::Value *L = LHS->codegen();
	llvm::Value *R = RHS->codegen();
	if (!L && !R)
		return nullptr;

	switch (Op) {
	case '+':
		return Builder->CreateFAdd(L, R, "addtmp");
	case '-':
		return Builder->CreateFSub(L, R, "subtmp");
	case '*':
		return Builder->CreateFMul(L, R, "multmp");
	case '<':
		L = Builder->CreateFCmpULT(L, R, "cmptmp");
		// Convert bool 0/1 to double 0.0 or 1.0
		return Builder->CreateUIToFP(L, llvm::Type::getDoubleTy(*TheContext), "booltmp");
	default:
		return LogErrorV("invalid binary operator");
	}
}

llvm::Value *CallExprAST::codegen() {
	llvm::Function *CalleeF = TheModule->getFunction(Callee);
	if (!CalleeF)
		return LogErrorV("Unknown function referenced");

	if (CalleeF->arg_size() != Args.size())
		return LogErrorV("Incorrect # arguments passed");

	std::vector<llvm::Value *> ArgsV;
	for (unsigned i = 0, e = Args.size(); i != e; ++i) {
		ArgsV.push_back(Args[i]->codegen());
		if (!ArgsV.back())
			return nullptr;
	}

	return Builder->CreateCall(CalleeF, ArgsV, "calltmp");
}

llvm::Function *PrototypeAST::codegen() {
	std::vector<llvm::Type*> Doubles(Args.size(), llvm::Type::getDoubleTy(*TheContext));
	llvm::FunctionType *FT = llvm::FunctionType::get(llvm::Type::getDoubleTy(*TheContext), Doubles, false);
	llvm::Function *F = llvm::Function::Create(FT, llvm::Function::ExternalLinkage, Name, TheModule.get());

	unsigned Idx = 0;
	for (auto &Arg : F->args())
		Arg.setName(Args[Idx++]);

	return F;
}

llvm::Function *FunctionAST::codegen() {
	llvm::Function *TheFunction = TheModule->getFunction(Proto->getName());

	if (!TheFunction)
		TheFunction = Proto->codegen();

	if (!TheFunction)
		return nullptr;

	if (!TheFunction->empty())
		return (llvm::Function*)LogErrorV("Function cannot be redefined");

	llvm::BasicBlock *BB = llvm::BasicBlock::Create(*TheContext, "entry", TheFunction);
	Builder->SetInsertPoint(BB);

	NamedValues.clear();
	for (auto &Arg : TheFunction->args())
		NamedValues[std::string(Arg.getName())] = &Arg;

	if (llvm::Value *RetVal = Body->codegen()) {
		Builder->CreateRet(RetVal);
		llvm::verifyFunction(*TheFunction);
		return TheFunction;
	}

	// TODO: Fix bug by validating the signature against the prototype
	// https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/LangImpl03.html#function-code-generation

	TheFunction->eraseFromParent();
	return nullptr;
}

int main() {
	BinopPrecedence['<'] = 10;
	BinopPrecedence['+'] = 20;
	BinopPrecedence['-'] = 30;
	BinopPrecedence['*'] = 40;

	fprintf(stderr, "ready> ");
	getNextToken();

	MainLoop();

	return 0;
}
