﻿using VehicleHandler.Actors;
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddActors(options =>
{
    // Register actor types and configure actor settings
    options.Actors.RegisterActor<VehicleActor>();
});

builder.Services.AddControllers();
builder.Services.AddDaprClient();
builder.Services.AddHttpClient();


var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}
else
{
    // By default, ASP.Net Core uses port 5000 for HTTP. The HTTP
    // redirection will interfere with the Dapr runtime. You can
    // move this out of the else block if you use port 5001 in this
    // example, and developer tooling (such as the VSCode extension).
    //app.UseHttpsRedirection();
}

app.MapActorsHandlers();
app.MapControllers();

app.UseCloudEvents();
app.MapSubscribeHandler();

app.Run();
