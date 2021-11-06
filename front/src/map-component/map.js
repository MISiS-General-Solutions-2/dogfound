import React, {useState, useEffect} from "react";
import ReactMapboxGl, {Layer, Feature, Marker, ZoomControl, Cluster} from 'react-mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import './map.css';

const Map = ReactMapboxGl({
    accessToken:
        'pk.eyJ1IjoiYmFkYmFkbm90Z29vZCIsImEiOiJja3RxMTdqdHkwcnRxMm5vYWVvcXVia3J5In0.97Gsy4fkvJQsrkD8_XeFLA',
});

let markerCoords = [];

export default function MapComponent(props) {
    const [data, setData] = useState(null);
    const scrollTo = props.scroll;
    useEffect(() => {
        if (props.data !== undefined) {
            setData(props.data.data);
        }
    }, [props.data]);

    return (
        <Map
            style="mapbox://styles/mapbox/streets-v9"
            containerStyle={{
                height: '100%',
                width: '100%'
            }}
            center={[props.lng, props.lat]}
        >
            <ZoomControl/>
            {data !== [] && data !== undefined && data !== null ?
                data.map((el, index) => (
                    el.lonlat[0] !== 0 && el.lonlat[1] !== 0 && el.timestamp !== 0 ?
                        <Marker coordinates={[el.lonlat[0], el.lonlat[1]]} anchor="bottom">
                            <div id={"marker_" + index}  className={'marker_div'} onClick={() => scrollTo(index)}/>
                        </Marker>
                        : null
                ))
                : null
            }
        </Map>

    );
}