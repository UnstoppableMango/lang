using Microsoft.AspNetCore.Server.Kestrel.Core;
using UnMango.Lang.Host.Services;

if (args.Length == 0) {
	Console.WriteLine("socket path is required");
	return 1;
}

var socketPath = Path.Join(AppContext.BaseDirectory, args[0]);

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

builder.WebHost.ConfigureKestrel(options => {
	options.ListenUnixSocket(socketPath, listen => {
		listen.Protocols = HttpProtocols.Http2;
	});
});

var app = builder.Build();

app.MapGrpcService<ParserService>();

app.Run();

return 0;
