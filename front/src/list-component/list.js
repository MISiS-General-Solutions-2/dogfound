import React from "react";
import './list.css';
const API_URL = process.env.REACT_APP_API_URL || "http://localhost:1022";

export default class ListComponent extends React.Component {
    constructor(props) {
        super(props);
    }
    render() {
        let dataTemp = this.props.data;
        const data = dataTemp.data;
        const action = this.props.action
        console.log(data);
        return (
            <div className="listContainer">
                <div className="listDiv">
                    {data !== null ?
                        data.map((el) => (
                            <button ket={el.filename} className="listButton">
                                {el.address !== '' ?
                                    <p className="listAddress">
                                        {el.address}
                                    </p>
                                    : null}
                                <p className="listAddress">
                                    {el.timestamp}
                                </p>
                                <img src={`http://${API_URL}/api/image/` + el.filename} alt="" />
                            </button>
                        ))
                        : null}
                </div>
                <div className="listResetDiv">
                    <button className="listReset" onClick={action}>
                        Попробовать еще раз
                    </button>
                </div>
            </div>
        )
    }
}