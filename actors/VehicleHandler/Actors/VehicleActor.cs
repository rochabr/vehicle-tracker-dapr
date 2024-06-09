using System.Threading;
using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;
using VehicleHandler.Models;


namespace VehicleHandler.Actors
{
    public interface IVehicleActor : IActor
    {
        Task<bool> StartShipment(Shipment shipmentId);
    }

    public class VehicleActor : Actor, IVehicleActor
    {

        public readonly string ShipmentStatusPending = "pending";
        public readonly string ShipmentStatusEnRoute = "en-route";
        public readonly string ShipmentStatusCompleted = "completed";
        string SHIPMENT_STATE_STORE = "vtd.shipment.state";
        string PUBSUB = "vtd.pubsub";
        string LOCATION_TOPIC = "locations";

        public VehicleActor(ActorHost host) : base(host)
        {
        }
        // public async Task SetState(string shipmentId)
        // {
        //     using var client = new DaprClientBuilder().Build();

        //     //Using Dapr SDK to save and get state
        //     await client.SaveStateAsync(STATESTORE, shipmentId, state);
        // }

        // public async Task<string> GetState()
        // {
        //     using var client = new DaprClientBuilder().Build();
        //     return await client.GetStateAsync<string>(STATESTORE, shipmentId);
        // }

        public async Task<bool> StartShipment(Shipment shipment)
        {
            using var client = new DaprClientBuilder().Build();

            if (shipment == null)
            {
                Console.WriteLine("Shipment not found");
                return false;
            }

            Console.WriteLine("Shipment path: " + shipment.Path.Positions.Count);

            foreach (Position position in shipment.Path.Positions)
            {
                Console.WriteLine("Position: " + position.Latitude + " - " + position.Longitude);

                // Publish each point in the path to the pub/sub
                var shipmentPosition = new ShipmentPosition(shipment.ShipmentId, position);

                //publish last position to pubsub
                try
                {
                    CancellationTokenSource source = new CancellationTokenSource();
                    CancellationToken cancellationToken = source.Token;

                    await client.PublishEventAsync(PUBSUB, LOCATION_TOPIC, shipmentPosition, cancellationToken);
                    Console.WriteLine("Published last position data: " + shipmentPosition);
                }
                catch (Exception e)
                {
                    Console.WriteLine("Error publishing shipment position for Shipment: {shipmentId}. Message: {Content}", shipment.ShipmentId, e.InnerException?.Message ?? e.Message);
                    return false;
                }

                Thread.Sleep(3000);
            }

            return true;

        }

        protected override Task OnActivateAsync()
        {
            // Provides an opportunity to perform some optional setup when an actor is activated.
            // An actor is activated the first time any of its methods are invoked.
            Console.WriteLine($"Shipment {this.Id} is taking off!");
            return Task.CompletedTask;
        }
    }
}