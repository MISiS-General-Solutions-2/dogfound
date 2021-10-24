import React from "react";
import './list.css';

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
                    {data.map((el) => (
                        <button ket={el.filename} className="listButton">
                            {el.address !== '' ?
                                <p className="listAddress">
                                    {el.address}
                                </p>
                                : null}
                            <p className="listAddress">
                                {el.timestamp}
                            </p>
                            <img src={"http://localhost:5000/api/image/" + el.filename} alt="" />
                        </button>
                    ))}
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