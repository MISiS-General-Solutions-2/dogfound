import React, {useState, useEffect} from "react";
import {Link, DirectLink, Element, Events, animateScroll as scroll, scrollSpy, scroller} from 'react-scroll'
import axios from "axios";
import Header from "../header-component/header";
import MapComponent from "../map-component/map";
import Search from "../search-component/search";
import './main.css';

export default function Main() {
    const [lat, setLat] = useState(55.75219808341882);
    const [id, setId] = useState('');
    const [lng, setLng] = useState(37.621985952011016);
    const [filename, setFilename] = useState(null);
    const [timestamp, setTimestamp] = useState(0);
    const [option1, setOption1] = useState(null);
    const [option2, setOption2] = useState(null);
    const [data, setData] = useState(undefined);
    const [imageShow, setImageShow] = useState(false);
    const [listShow, setListShow] = useState(false);
    const address = window.location.href;
    console.log(lat, lng);

    function scrollTo(e) {
        if (id !== '' && id !== null) {
            let elems = document.getElementsByClassName("listButton");
            [].forEach.call(elems, function (el) {
                el.classList.remove("Focus");
            });
            document.getElementById(id).classList.remove("Focus");
            console.log(document.getElementsByClassName('Focus'));
        }

        setId(e);
        document.getElementById(e).classList.add('Focus');
        scroller.scrollTo('img_' + e, {
            duration: 500,
            smooth: true,
            containerId: 'list_container',
        })
        let elm = document.getElementsByClassName("listButton")
    }

    useEffect(() => {
        Events.scrollEvent.register('begin', function () {
            console.log("begin", arguments);
        });

        Events.scrollEvent.register('end', function () {
            console.log("end", arguments);
        });
    })

    function sendData(e1, e2, e3, e4) {
        let dataTemp = {};
        dataTemp["t0"] = e1;
        dataTemp["t1"] = e2;
        if (e3 !== 0) {
            dataTemp["color"] = e3;
        }
        if (e4 !== 0) {
            dataTemp["tail"] = e4;
        }
        const url = address + 'api/image/by-classes';
        // const url = 'http://5.228.244.67:1022/api/image/by-classes';
        console.log(url);
        const options = {
            method: 'POST',
            headers: {'content-type': 'application/json'},
            data: dataTemp,
            url
        };
        axios(options)
            .then(response => {
                setData(response);

                console.log(e1, e2, e3, e4);
            });
    }

    return (
        <main-screen>
            <Header/>
            <main-component id="mainapp">
                <Search action1={setTimestamp} action2={setOption1} action3={setOption2} action4={sendData}
                        action5={setData} list={setListShow} listStatus={listShow} data={data} setLat={setLat}
                        setLng={setLng}/>
                <MapComponent scroll={scrollTo} lat={lat} lng={lng} data={data} openStatus={imageShow}
                              show={setImageShow} setId={setId}/>
            </main-component>
        </main-screen>
    )
}