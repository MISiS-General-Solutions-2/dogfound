export default function arr() {
    let arr = ['1633392000',
        '1632787200',
        '1632528000',
        '1632700800',
        '1633305600',
        '1632873600',
        '1633478400',
        '1632441600',
        '1633046400',
    ]
    let reArr = [];
    arr.map((el) => {
        let dt = new Date(el * 1000);
        let tempObj = { day: dt.getUTCDate(), month: dt.getMonth(), year: dt.getFullYear() }
        reArr.push(tempObj);
    });
    reArr.map((el) => {
        console.log(el.day + ':' + el.month + ':' + el.year);
    })
    return reArr;
}

