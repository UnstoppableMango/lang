using Grpc.Core;

namespace UnMango.Lang.Host.Services;

using static Lang.ParserService;

public sealed class ParserService : ParserServiceBase
{
	public override Task<ParseResponse> Parse(ParseRequest request, ServerCallContext context)
	{
		var result = Parser.parse(request.Text);

		if (!result.IsError)
			return Task.FromResult(new ParseResponse
			{
				File = ToFile(result.ResultValue),
			});

		return Task.FromException<ParseResponse>(new Exception(result.ErrorValue));
	}

	private static File ToFile(Ast.File file)
	{
		File f = new();
		f.Nodes.Add(file.Nodes.Select(ToNode));
		return f;
	}

	private static Node ToNode(Ast.Node node) =>
		new()
		{
			String =
			{
				Value = node.Item,
			},
		};
}
