using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace VehicleHandler.Models
{
    public class Vehicle
    {
        [JsonPropertyName("vehicleId")]
        public int VehicleId { get; set; }

        [JsonPropertyName("driver")]
        public required string Driver { get; set; }
    }
}
