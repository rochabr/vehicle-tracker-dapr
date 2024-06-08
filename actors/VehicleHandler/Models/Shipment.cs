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
        public Path path { get; set; }

        [JsonPropertyName("status")]
        public string Status { get; set; }
    }
    public class Path
    {
        [JsonPropertyName("positions")]
        public List<Position> positions { get; set; }
    }

    public class Position
    {
        [JsonPropertyName("_lat")]
        public string Latitude { get; set; }

        [JsonPropertyName("_lon")]
        public string Longitude { get; set; }
    }

    public class ShipmentPosition
    {
        [JsonPropertyName("shipmentId")]
        public string ShipmentId { get; set; }

        [JsonPropertyName("position")]
        public Position Position { get; set; }
    }
}

