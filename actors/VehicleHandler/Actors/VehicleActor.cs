using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;
using VehicleHandler.Models;


namespace VehicleHandler.Actors
{
    public interface IVehicleActor : IActor
    {
        // Task<string> GetState();
        // Task SetState(string state);
        Task<string> SayHelloWorld();
        Task<bool> DeleteOrder(string orderId);
    }

    public class VehicleActor : Actor, IVehicleActor
    {
        string STATESTORE = "vtd.shipment.state";
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

        public Task<string> SayHelloWorld()
        {
            return Task.FromResult("Hello World!");
        }

        public async Task<bool> DeleteOrder(string orderId)
        {
            try
            {
                using var client = new DaprClientBuilder().Build();

                CancellationTokenSource source = new CancellationTokenSource();
                CancellationToken cancellationToken = source.Token;

                await client.DeleteStateAsync(STATESTORE, orderId, cancellationToken: cancellationToken);

                return true;
            }
            catch (Exception ex)
            {
                return false;
            }

        }

        private async Task<Shipment> GetShipment(string shipmentId)
        {
            using var client = new DaprClientBuilder().Build();
            return await client.GetStateAsync<Shipment>(STATESTORE, shipmentId);
        }
    }
}