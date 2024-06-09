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
        [HttpPost("start-shipment")]
        public IActionResult StartVehicle([FromBody] Shipment shipment)
        {
            var vehicleActorProxy = ActorProxy.Create<IVehicleActor>(new ActorId(shipment.ShipmentId), VehicleActorType);
            vehicleActorProxy.StartShipment(shipment);

            return Ok();

        }
    }
}
