import React from "react";
import MapGL, { Marker } from 'react-map-gl';
import ReactModal from 'react-modal';
import ModalComponent from "../modal-component/modal";
import 'mapbox-gl/dist/mapbox-gl.css';
import './map.css';

// mapboxgl.accessToken = 'pk.eyJ1IjoiYmFkYmFkbm90Z29vZCIsImEiOiJja3RxMTdqdHkwcnRxMm5vYWVvcXVia3J5In0.97Gsy4fkvJQsrkD8_XeFLA';

const MAPBOX_TOKEN = 'pk.eyJ1IjoiYmFkYmFkbm90Z29vZCIsImEiOiJja3RxMTdqdHkwcnRxMm5vYWVvcXVia3J5In0.97Gsy4fkvJQsrkD8_XeFLA'


export default class MapComponent extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            lng: null,
            lat: null,
            center: [37.621985952011016,],
            viewport: {
                latitude: 55.75219808341882,
                longitude: 37.621985952011016,
                zoom: 12,
                bearing: 0,
                pitch: 0
            }
        }
        this.showModal = this.showModal.bind(this);
    }
    showModal(filename) {
        console.log(filename)
        const show = this.props.show;
        this.setState({
            filename: filename
        })
        show();
    }

    // componentDidMount() {
    //     const map = new mapboxgl.Map({
    //         container: 'map', // container ID
    //         style: 'mapbox://styles/mapbox/outdoors-v11', // style URL
    //         center: [37.621985952011016, 55.75219808341882], // starting position [lng, lat]
    //         zoom: 12 // starting zoom
    //     });
    //     map.addControl(new mapboxgl.NavigationControl());
    //     map.addControl(
    //         new mapboxgl.GeolocateControl({
    //             positionOptions: {
    //                 enableHighAccuracy: true
    //             },
    //             // When active the map will receive updates to the device's location as it changes.
    //             trackUserLocation: true,
    //             // Draw an arrow next to the location dot to indicate which direction the device is heading.
    //             showUserHeading: true
    //         })
    //     );
    //     const language = new MapboxLanguage();
    //     map.addControl(language);
    // }
    render() {
        ReactModal.setAppElement('#root');
        let dataTemp = this.props.data;
        if (dataTemp === null) {
            dataTemp = [];
        }
        const hide = this.props.hide
        const openStatus = this.props.openStatus;
        const data = dataTemp.data;
        const { filename } = this.state;
        console.log(data);
        return (
            // <div id={"map"} />
            <MapGL
                id={"map"}
                {...this.state.viewport}
                width="100%"
                height="100%"
                mapStyle="mapbox://styles/mapbox/outdoors-v11"
                onViewportChange={viewport => this.setState({ viewport })}
                mapboxApiAccessToken={MAPBOX_TOKEN}
            >
                {data !== undefined && data !== null && data !== [] ?
                    data.map((el) => (
                        el.lonlat[0] !== 0 && el.lonlat[1] !== 0 ?
                            <Marker key={el.filename} latitude={el.lonlat[0]} longitude={el.lonlat[1]} offsetLeft={0} offsetTop={-40} onClick={() => this.showModal(el.filename)}>
                                <div className="markerDiv">
                                    {/* <img className={"imgMarker"} src={"http://localhost:1022/api/image/" + el.filename} alt="" /> */}
                                </div>
                            </Marker>
                            : null
                    ))
                    : null}
                <ReactModal
                    isOpen={openStatus}
                >
                    <ModalComponent action={hide} data={filename} />
                </ReactModal>
            </MapGL>
        );
    }
}
