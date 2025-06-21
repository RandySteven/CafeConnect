import {useEffect, useRef} from "react";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

export interface Coordinate  {
    long: number;
    lat: number;
};


export const useGetLocation = () => {
    return
}

export const useAWSLocationServiceMultiPinPoint = (coordinates: Coordinate[]) => {
    const mapContainer = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!mapContainer.current || coordinates.length === 0) return;

        const region = "us-east-1";
        const style = "Standard";
        const apiKey = "v1.public.eyJqdGkiOiJmMTVhNTFhYi1jZWJiLTQzN2EtYjZlYi1mNTEwMzk0ZTdlZTMifVgAo9b8xNkOfU8EbbvGJ-PsmsLfLR9oi_wY29CYpmulQUsLdCcsuK5IRhFeOQUlP4m_hzjVYdcSWd762KFxXVVMlSXBO20AcJs0E59rkRdoalC8vwUSBCP4yAwF7hwjwzg4BeMLW7p4qefJGi0mE4jxVqVZ-IG_suK7A33ZqN7yXaQcgr6JH-xc--ghStXvTRpD9X0rR6YQDDBd6JDFZIGkcQfYZrLNU6Un8Fi8K6ccO_kudKGZkTMhl5LnFdiJMhlkgr1QxiCJ3IIsXtdGUdbMnVSjjC62wmSqzeWLnYOjkHYzCAlDmAzq25GYluyes9RBr7Ho4BcPFE90-tYX1QU.ZWU0ZWIzMTktMWRhNi00Mzg0LTllMzYtNzlmMDU3MjRmYTkx";
        const colorScheme = "Light";

        // Use the first coordinate to center the map
        let { long, lat } = coordinates[0];

        const map = new maplibregl.Map({
            container: mapContainer.current,
            style: `https://maps.geo.${region}.amazonaws.com/v2/styles/${style}/descriptor?key=${apiKey}&color-scheme=${colorScheme}`,
            center: [long, lat],
            zoom: 12,
        });

        map.addControl(new maplibregl.NavigationControl(), "top-left");

        // Add multiple markers
        coordinates.forEach(({ long, lat }) => {
            new maplibregl.Marker({ color: "red" })
                .setLngLat([long, lat])
                .addTo(map);
        });

        return () => map.remove();
    }, [coordinates]);

    return mapContainer;
};

export const useAWSLocationServiceSinglePinPoint = (coordinate :Coordinate) => {
    if(isNaN(coordinate.long) || isNaN(coordinate.lat)) {
        coordinate.long = 0
        coordinate.lat = 0
    }
    const mapContainer = useRef<HTMLDivElement>(null)
    useEffect(() => {
        const region = "us-east-1"; // Your region
        const style = "Standard"
        const apiKey = "v1.public.eyJqdGkiOiJmMTVhNTFhYi1jZWJiLTQzN2EtYjZlYi1mNTEwMzk0ZTdlZTMifVgAo9b8xNkOfU8EbbvGJ-PsmsLfLR9oi_wY29CYpmulQUsLdCcsuK5IRhFeOQUlP4m_hzjVYdcSWd762KFxXVVMlSXBO20AcJs0E59rkRdoalC8vwUSBCP4yAwF7hwjwzg4BeMLW7p4qefJGi0mE4jxVqVZ-IG_suK7A33ZqN7yXaQcgr6JH-xc--ghStXvTRpD9X0rR6YQDDBd6JDFZIGkcQfYZrLNU6Un8Fi8K6ccO_kudKGZkTMhl5LnFdiJMhlkgr1QxiCJ3IIsXtdGUdbMnVSjjC62wmSqzeWLnYOjkHYzCAlDmAzq25GYluyes9RBr7Ho4BcPFE90-tYX1QU.ZWU0ZWIzMTktMWRhNi00Mzg0LTllMzYtNzlmMDU3MjRmYTkx"
        const colorScheme = "Light";

        const map = new maplibregl.Map({
            container: mapContainer.current!,
            style: `https://maps.geo.${region}.amazonaws.com/v2/styles/${style}/descriptor?key=${apiKey}&color-scheme=${colorScheme}`,
            center: [coordinate.long, coordinate.lat],
            zoom: 11,
    });
        map.addControl(new maplibregl.NavigationControl(), "top-left")
        new maplibregl.Marker({ color: "red" })
            .setLngLat([coordinate.long, coordinate.lat])
            .addTo(map)
    }, [coordinate.long, coordinate.lat]);

    return mapContainer
}