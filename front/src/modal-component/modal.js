import React from "react";
import './modal.css'
const API_URL = process.env.REACT_APP_API_URL || "http://localhost:1022";

export default class ModalComponent extends React.Component {
    constructor(props) {
        super(props);
    }
    render() {
        const action = this.props.action;
        const filename = this.props.data;
        return (
            <div className="modalWindow">
                <button className="closeButton" onClick={action}></button>
                <img src={`http://${API_URL}/api/image/` + filename} alt="" />
            </div>
        )
    }
}