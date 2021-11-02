import React from "react";
import './modal.css'

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
                <img src={"http://5.228.244.67:1022/api/image/" + filename} alt="" />
            </div>
        )
    }
}