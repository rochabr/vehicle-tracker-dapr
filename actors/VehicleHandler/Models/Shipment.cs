using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace VehicleHandler.Models
{

    public class Shipment
    {
        [JsonPropertyName("shipmentId")]
        public string ShipmentId { get; set; }

        [JsonPropertyName("vehicle")]
        public Vehicle Vehicle { get; set; }

        [JsonPropertyName("path")]
        public Path Path { get; set; }

        [JsonPropertyName("status")]
        public string Status { get; set; }
    }
}


