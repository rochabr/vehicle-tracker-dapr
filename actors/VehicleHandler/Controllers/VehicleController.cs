//dependencies 
using System.Collections.Generic;
using System.Threading.Tasks;
using System;
using Microsoft.AspNetCore.Mvc;
using Dapr;
using Dapr.Client;
using VehicleHandler.Models;
using VehicleHandler.Actors;
using Dapr.Actors;
using Dapr.Actors.Client;

//code
namespace VehicleHandler.Controllers
{
    [ApiController]
    public class VehicleController : Controller
    {

        private const string VehicleActorType = "VehicleActor";

        //Subscribe to a topic 
        [Topic("vtd.pubsub", "shipments")]
        [HttpPost("start-shipment")]
        public async Task<string> startVehicle([FromBody] Shipment shipment)
        {
            Console.WriteLine("Subscriber received : " + shipment.ShipmentId);

            var vehicleActorProxy = ActorProxy.Create<IVehicleActor>(new ActorId(shipment.ShipmentId), VehicleActorType);
            var result = await vehicleActorProxy.SayHelloWorld();
            Console.WriteLine("Subscriber received : " + result);

            return result;
        }
    }
}
