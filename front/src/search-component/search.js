import React, { useState } from "react";
import { DateRange, Calendar } from 'react-date-range';
import Select from 'react-select'
import Modal from 'react-modal';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import './array';
import './search.css';
import './array';
import ListComponent from "../list-component/list";
import axios from "axios";
import arr from "./array";

Modal.setAppElement('body');

export default function Search(props) {

    const [time0, setTime0] = useState(null);
    const [time1, setTime1] = useState(null);
    const [value1, setValue1] = useState(null);
    const [value2, setValue2] = useState(null);
    const [error1, setError1] = useState(false);
    const [error2, setError2] = useState(false);
    const [error3, setError3] = useState(false);

    const [showStatus, setShowStatus] = useState(false);
    const [date, setDate] = useState(null);
    const [imgValue1, setImageValue1] = useState(null);
    const [imgValue2, setImageValue2] = useState(null);

    const options1 = [
        { value: 0, label: 'Любой' },
        { value: 1, label: 'Cветлая' },
        { value: 2, label: 'Темная' },
        { value: 3, label: 'Разноцветная' }
    ];
    const options2 = [
        { value: 0, label: 'Любой' },
        { value: 1, label: 'Длинный хвост' },
        { value: 2, label: 'Короткий хвост' }
    ];

    const data = props.data;
    const action4 = props.action4
    const action5 = props.action5;

    const setLat = props.setLat;
    const setLng = props.setLng;

    function sendDataButton() {
        setError1(time0 === null || time1 === null ? true : false);
        setError2(value1 === null ? true : false);
        setError3(value2 === null ? true : false);
        if (error1 !== true && error2 !== true && error3 !== true && time0 !== null && time1 !== null && value1 !== null && value2 !== null) {
            console.log('sss');
            action4(time0.getTime() / 1000, time1.getTime() / 1000, value1.value, value2.value)
        }
    }

    const [state, setState] = useState([
        {
            startDate: new Date(),
            endDate: new Date(),
            key: 'selection'
        }
    ]);

    if (data === undefined) {
        return (
            <search-component id={"search-component"}>
                <p>Поиск питомца</p>
                <div className="option">Дата пропажи
                    {error1 ? <div className="error">Поле не заполнено!</div> : null}
                </div>
                <div className={'date_holder'}>
                    <DateRange
                        editableDateInputs={true}
                        onChange={item => {
                            setState([item.selection]);
                            setTime0(item.selection.startDate);
                            setTime1(item.selection.endDate);
                            setError1(false);
                            console.log('MOVE');
                        }}
                        moveRangeOnFirstSelection={false}
                        ranges={state}
                    />
                </div>
                <div className="option">Окрас
                    {error2 ? <div className="error">Поле не заполнено!</div> : null}
                </div>
                <Select options={options1} defaultValue={0} placeholder={"Цвет"} onChange={(e) => {
                    setValue1(e);
                    setError2(false);
                }} />
                <div className="option">Тип хвоста
                    {error3 ? <div className="error">Поле не заполнено!</div> : null}
                </div>
                <Select options={options2} defaultValue={0} placeholder={"Хвост"} onChange={(e) => {
                    setValue2(e);
                    setError3(false);
                }} />
                <button className={"SendButton"} onClick={sendDataButton}>Найти!</button>
                <button className={"addImage"} onClick={() => setShowStatus(true)}>
                    Добавить фото
                </button>
                <Modal
                    isOpen={showStatus}
                    onRequestClose={() => setShowStatus(false)}
                >
                    <p>Помощь в поиске собак</p>
                    <p>Фотография питомца</p>
                    <input type="file" id="fileInput" />
                    <p>Дата пропажи</p>
                    <Calendar
                        date={date}
                        onChange={item => setDate(item)}
                    />
                    <input type={"number"} id={"lonImg"} placeholder={"Долгота"} />
                    <input type={"number"} id={"latImg"} placeholder={"Широта"} />
                    <button className={"sendImgFinal"} onClick={() => {
                        let filedata = document.getElementById("fileInput");
                        let latImg = document.getElementById('lonImg').value;
                        let lonImg = document.getElementById('latImg').value;
                        if (filedata.files.length !== 0 && lonImg !== (null || '') && latImg !== (null || '')) {
                            console.log('da');
                            arr();
                            let formData = new FormData();
                            // let url = window.location.href + 'api/image/upload';
                            let url = window.location.href + "api/image/upload?timestamp=" + date.getTime() / 1000 + '&lon=' + lonImg + '&lat=' + latImg;
                            formData.append('file', filedata.files[0]);
                            const options = {
                                method: 'PUT',
                                headers: { 'content-type': 'multipart/form-data' },
                                data: formData,
                                url
                            };
                            axios(options)
                                .then(response => {
                                    console.log(response);
                                });
                            setShowStatus(false)
                        } else {
                            alert('Введите данные!');
                        }
                    }}>
                        Отправить
                    </button>
                    <p className={"addP"}>Помогите найти потерявшуюся собаку! <br /> Если вы обнаружили на улице животное
                        без хозяина, <br /> сфотографируйте его как можно лучше и отправьте нам, <br /> а мы постараемся
                        помочь хозяину его найти.</p>
                </Modal>
            </search-component>
        );
    } else {
        return (
            <ListComponent data={data} action1={action5} setLat={setLat} setLng={setLng} />
        )
    }
}