using System.Text.Json.Serialization;
namespace VehicleHandler.Models
{
    public class Path
    {
        [JsonPropertyName("positions")]
        public List<Position>? Positions { get; set; }
    }
}