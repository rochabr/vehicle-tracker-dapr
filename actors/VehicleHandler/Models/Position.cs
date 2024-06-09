using System.Text.Json.Serialization;
namespace VehicleHandler.Models
{
    public class Position
    {
        [JsonPropertyName("_lat")]
        public double Latitude { get; set; }

        [JsonPropertyName("_lon")]
        public double Longitude { get; set; }
    }
}
