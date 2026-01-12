using Grpc.Core;

namespace grpcs.Services;

public class GreeterService : Greeter.GreeterBase
{
	private readonly ILogger<GreeterService> logger;
	public GreeterService(ILogger<GreeterService> logger)
	{
		this.logger = logger;
	}

	public override Task<HelloReply> SayHello(
		HelloRequest req,
		ServerCallContext ctx
		)
	{
		logger.LogInformation(
			"SayHello: name={name} peer={peer}",
			req.Name,
			ctx.Peer
			);
		return Task.FromResult(new HelloReply
		{
			Message = "Hello " + req.Name
		});
	}
}
