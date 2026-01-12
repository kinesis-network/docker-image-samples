using Grpc.Net.Client;
using grpcs;
using grpcs.Services;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using System.CommandLine;

var rootCommand = new RootCommand
{
	Description = "gRPC server/client executable",
};

var serverCommand = new Command("server")
{
	Description = "Run gRPC server"
};
serverCommand.SetAction(_ =>
{
	RunServer();
});

rootCommand.Subcommands.Add(serverCommand);

var addressOption = new Option<string>("--address")
{
	Description = "gRPC server address",
	DefaultValueFactory = _ => "http://localhost:5000",
};
var nameOption = new Option<string>("--name")
{
	Description = "Name to send to server",
	DefaultValueFactory = _ => ":)",
};
var clientCommand = new Command("client")
{
	Description = "Run gRPC client",
};
clientCommand.Options.Add(addressOption);
clientCommand.Options.Add(nameOption);
clientCommand.SetAction(async (ParseResult parseResult, CancellationToken ct) =>
{
	var address = parseResult.GetValue(addressOption);
	var name = parseResult.GetValue(nameOption);
	await RunClientAsync(address!, name!);
});

rootCommand.Subcommands.Add(clientCommand);

var parseResult = rootCommand.Parse(args);
return await parseResult.InvokeAsync();

static void RunServer()
{
	var builder = WebApplication.CreateBuilder();
	builder.WebHost.ConfigureKestrel(options =>
	{
		// This tells Kestrel: "Whatever port you end up listening on
		// (from launchSettings, env vars, etc.), use HTTP/2."
		options.ConfigureEndpointDefaults(listenOptions =>
		{
			listenOptions.Protocols = HttpProtocols.Http2;
		});
	});
	builder.Services.Configure<KestrelServerOptions>(options =>
	{
		options.AllowAlternateSchemes = true;
	});
	builder.Services.AddGrpc();

	var app = builder.Build();
	app.MapGrpcService<GreeterService>();
	app.MapGet("/", () =>
		"Communication with gRPC endpoints must be made through a gRPC client");

	app.Run();
}

static async Task RunClientAsync(string address, string name)
{
	using var channel = GrpcChannel.ForAddress(address);
	var client = new Greeter.GreeterClient(channel);
	var reply = await client.SayHelloAsync(new HelloRequest
	{
		Name = name
	});
	Console.WriteLine($"Server replied: {reply.Message}");
}
