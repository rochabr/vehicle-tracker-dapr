using System.Text.Json.Serialization;

namespace VehicleHandler.Models
{
    public class ShipmentPosition
    {
        public ShipmentPosition(string shipmentId, Position position)
        {
            ShipmentId = shipmentId;
            Position = position;
        }

        [JsonPropertyName("shipmentId")]
        public string ShipmentId { get; set; }

        [JsonPropertyName("position")]
        public Position Position { get; set; }
    }
}
