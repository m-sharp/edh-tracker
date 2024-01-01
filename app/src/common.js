// Record takes a record dictionary like {1: 10, 2: 12, 3: 7, 4: 5}
export function Record({ record }) {
    let first = record[1] || 0;
    let second = record[2] || 0;
    let third = record[3] || 0;
    let fourth = record[4] || 0;

    return (
        <span id="record">{first} / {second} / {third} / {fourth}</span>
    )
}
