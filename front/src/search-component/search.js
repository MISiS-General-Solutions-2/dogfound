import React, {useState} from "react";
import {DateRange} from 'react-date-range';
import Select from 'react-select'
import 'react-date-range/dist/styles.css'; // main css file
import 'react-date-range/dist/theme/default.css'; // theme css file
import './search.css';
import ListComponent from "../list-component/list";

export default function Search(props) {
    const [time0, setTime0] = useState(null);
    const [time1, setTime1] = useState(null);
    const [value1, setValue1] = useState(null);
    const [value2, setValue2] = useState(null);
    const [error1, setError1] = useState(false);
    const [error2, setError2] = useState(false);
    const [error3, setError3] = useState(false);

    const options1 = [
        {value: 0, label: 'Любой'},
        {value: 1, label: 'Cветлая'},
        {value: 2, label: 'Темная'},
        {value: 3, label: 'Разноцветная'}
    ];
    const options2 = [
        {value: 0, label: 'Любой'},
        {value: 1, label: 'Длинный хвост'},
        {value: 2, label: 'Короткий хвост'}
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
            <search-component>
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
                }}/>
                <div className="option">Тип хвоста
                    {error3 ? <div className="error">Поле не заполнено!</div> : null}
                </div>
                <Select options={options2} defaultValue={0} placeholder={"Хвост"} onChange={(e) => {
                    setValue2(e);
                    setError3(false);
                }}/>
                <button className={"SendButton"} onClick={sendDataButton}>Найти!</button>
            </search-component>
        );
    } else {
        return (
            <ListComponent data={data} action1={action5} setLat={setLat} setLng={setLng}/>
        )
    }
}