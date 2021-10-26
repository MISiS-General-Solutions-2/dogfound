import React from "react";
import DatePicker from 'react-date-picker';
import Select from 'react-select'
import './search.css';
import ListComponent from "../list-component/list";

export default class Search extends React.Component {
    constructor() {
        super();
        this.state = {
            value: "",
            error1: false,
            error2: false,
            error3: false,
            value1: null,
            value2: null,
            value3: null,
        }
        this.dateChange = this.dateChange.bind(this);
        this.select1 = this.select1.bind(this);
        this.select2 = this.select2.bind(this);
        this.sendDataButton = this.sendDataButton.bind(this);
    }
    dateChange(e) {
        const action1 = this.props.action1;
        this.setState({
            value: e
        })
        action1(e.getTime() / 1000);
        console.log(e.getTime());
        console.log(new Date().getTime());
        this.setState({
            error1: false,
            value1: e.getTime() / 1000,
        })
    }
    select1(e) {
        const action2 = this.props.action2;
        action2(e.value);
        this.setState({
            error2: false,
            value2: e.value,
        })
    }
    select2(e) {
        const action3 = this.props.action3;
        action3(e.value);
        console.log(e.value);
        this.setState({
            error3: false,
            value3: e.value,
        })
    }
    sendDataButton() {
        const { value1, value2, value3 } = this.state;
        const action4 = this.props.action4;
        this.setState({
            error1: value1 === null ? true : false,
            error2: value2 === null ? true : false,
            error3: value3 === null ? true : false,
        })
        if (value1 && value2 && value3 !== null) {
            action4();
        }
    }
    render() {
        const options1 = [
            { value: 1, label: 'Темная' },
            { value: 2, label: 'Cветлая' },
            { value: 3, label: 'Разноцветная' }
        ]
        const options2 = [
            { value: 1, label: 'Короткий хвост' },
            { value: 2, label: 'Длинный хвост' }
        ]
        const { value, error1, error2, error3 } = this.state;
        const data = this.props.data;
        const action = this.props.action;
        const action5 = this.props.action5;
        console.log(data);
        if (data === null) {
            return (
                <search-component>
                    <p>Поиск питомца</p>
                    <div className="option">Дата пропажи
                        {error1 ? <div className="error">Поле не заполнено!</div> : null}
                    </div>
                    <DatePicker
                        onChange={this.dateChange}
                        value={value}
                        defaultValue=""
                        dateFormat="MM/DD/YYYY"
                    />
                    <div className="option">Окрас
                        {error2 ? <div className="error">Поле не заполнено!</div> : null}
                    </div>
                    <Select options={options1} placeholder={"Цвет"} onChange={(e) => this.select1(e)} />
                    <div className="option">Тип хвоста
                        {error3 ? <div className="error">Поле не заполнено!</div> : null}
                    </div>
                    <Select options={options2} placeholder={"Хвост"} onChange={(e) => this.select2(e)} />
                    <button className={"SendButton"} onClick={this.sendDataButton}>Найти!</button>
                </search-component>
            );
        } else {
            return (
                <ListComponent data={data} action={action} action2={action5} />
            )
        }
    }
}