from runtime import Args
from typings.parking_lots_select.parking_lots_select import Input, Output
import requests

def handler(args: Args[Input])->Output:
    address = args.input.address
    response = requests.get(f"https://restapi.amap.com/v3/geocode/geo?key=your_api_key&address={address}&output=json")
    data = response.json()
    if "geocodes" not in data:
        raise ValueError(f"Failed to convert address to coordinates: {data}")
    geocode = data["geocodes"][0]
    location = geocode["location"]
    return {"longitude": float(location.split(",")[0]), "latitude": float(location.split(",")[1])}

