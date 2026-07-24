// Live Telemetry Polling for 7 Cars
function pollFleetLocations() {
    fetch('/api/fleet')
        .then(res => res.json())
        .then(data => {
            Object.keys(data).forEach(carId => {
                const vehicle = data[carId];
                
                // 1. Create or update directional marker
                if (!markers[carId]) {
                    const carIcon = L.divIcon({
                        className: 'custom-car-icon',
                        html: `<div style="transform: rotate(${vehicle.heading}deg); font-size: 24px;">🏎️</div>`,
                        iconSize: [30, 30]
                    });
                    
                    markers[carId] = L.marker([vehicle.lat, vehicle.lng], { icon: carIcon }).addTo(map)
                        .bindPopup(`<b>${vehicle.car_id}</b><br>Speed: ${Math.round(vehicle.speed * 3.6)} km/h<br>Heading: ${Math.round(vehicle.heading)}°`);
                    
                    routes[carId] = L.polyline([[vehicle.lat, vehicle.lng]], { color: '#38bdf8', weight: 4 }).addTo(map);
                } else {
                    // Update location and rotation angle
                    markers[carId].setLatLng([vehicle.lat, vehicle.lng]);
                    const iconDiv = markers[carId].getElement().querySelector('div');
                    if (iconDiv) {
                        iconDiv.style.transform = `rotate(${vehicle.heading}deg)`;
                    }
                    // Append coordinate to trailing route line
                    routes[carId].addLatLng([vehicle.lat, vehicle.lng]);
                }
            });
        })
        .catch(err => console.error("Error polling fleet data:", err));
}

// Poll backend every 2 seconds
setInterval(pollFleetLocations, 2000);
