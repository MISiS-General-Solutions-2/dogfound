import React, { useEffect, useState } from "react";
import './list.css';

export default function ListComponent(props) {
    const data = props.data.data;
    const action = props.action1;
    const [markerId, setMarkerId] = useState('');
    const setLat = props.setLat;
    const setLng = props.setLng;

    let options = {
        year: 'long',
        month: 'long',
        day: 'long',
        timezone: 'UTC'
    };

    let map = props.map;

    console.log(data);
    return (
        <div className="listContainer">
            <div className="listDiv" id={'list_container'}>
                {data !== [] && data !== undefined && data !== null ? data.map((el, index) => (
                    <button name={'img_' + index} id={index} key={el.filename} className="listButton"
                        onClick={() => {
                            if (el.timestamp !== 0 && el.lonlat[0] !== 0 && el.lonlat[1] !== 0) {
                                setLng(el.lonlat[0]);
                                setLat(el.lonlat[1]);
                                if (markerId !== '' && markerId !== null) {
                                    document.getElementById('marker_' + markerId).classList.remove("FocusMarker");
                                }
                                setMarkerId(index);
                                document.getElementById('marker_' + index).classList.add('FocusMarker');
                            }
                        }}
                    >
                        {el.address !== '' ?
                            <p className="listAddress">
                                {el.address}
                            </p>
                            : null}
                        {el.timestamp !== 0 ?
                            <p className="listAddress">
                                {new Date(el.timestamp * 1000).toLocaleDateString()}
                            </p>
                            : null}
                        {el.breed !== "" ?
                            <p>
                                {el.breed}
                            </p>
                            : null}
                        <img src={"http://5.228.244.67:1022/api/image/" + el.filename} alt="" />
                    </button>
                )) : null}
            </div>
            <div className="listResetDiv">
                <button className="listReset" onClick={() => action(undefined)}>
                    Попробовать еще раз
                </button>
            </div>
        </div>
    )
}
