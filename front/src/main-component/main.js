import React from "react";
import axios from "axios";
import Header from "../header-component/header";
import MapComponent from "../map-component/map";
import Search from "../search-component/search";
import './main.css';
const API_URL = process.env.REACT_APP_API_URL || "http://localhost:1022";

class Main extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            filename: null,
            timestamp: 0,
            option1: null,
            option2: null,
            data: null,
            imageShow: false
        }
        this.timestampGet = this.timestampGet.bind(this);
        this.option1Get = this.option1Get.bind(this);
        this.option2Get = this.option2Get.bind(this);
        this.sendData = this.sendData.bind(this);
        this.handleOpenModal = this.handleOpenModal.bind(this);
        this.handleCloseModal = this.handleCloseModal.bind(this);
        this.resetData = this.resetData.bind(this);
    }

    timestampGet(e) {
        this.setState({
            timestamp: e
        })
        console.log(e);
    }

    option1Get(e) {
        this.setState({
            option1: e
        })
        console.log(e);
    }

    option2Get(e) {
        this.setState({
            option2: e
        })
        console.log(e);
    }

    handleOpenModal() {
        this.setState({ imageShow: true });
    }

    handleCloseModal() {
        this.setState({ imageShow: false });
    }

    resetData() {
        console.log('reseted');
        this.setState({
            data: null,
            imageShow: false
        })
    }

    sendData() {
        const { option1, option2, timestamp } = this.state;
        let data = {
            "color": option1,
            "tail": option2,
            "timestamp": timestamp,
        }
        const url = `http://${API_URL}/api/image/by-classes`
        const options = {
            method: 'POST',
            headers: { 'content-type': 'application/json' },
            data: data,
            url,
        };
        axios(options)
            .then(response => {
                this.setState({
                    data: response
                })
            });
    }

    render() {
        const { data, imageShow, listShow } = this.state;
        return (
            <main-screen>
                <Header />
                <main-component id="mainapp">
                    <Search action1={this.timestampGet} action2={this.option1Get} action3={this.option2Get} action4={this.sendData} show={this.handleOpenModalList} hide={this.handleCloseModalList} status={listShow} data={data} action={this.resetData} />
                    <MapComponent data={data} openStatus={imageShow} show={this.handleOpenModal} hide={this.handleCloseModal} />
                </main-component>
            </main-screen>
        )
    }
}

export default Main;