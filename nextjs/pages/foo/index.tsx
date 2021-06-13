import styles from '../../styles/index.module.sass'
import Link from "next/link";
import useSWR from "swr";
import React from 'react';

async function fetcher(url: string) {
    const resp = await fetch(url);
    return resp.text();
}

const fetcherJson = url => fetch(url).then(r => r.json())

function Foo(): JSX.Element {

    const registerPeople = async event => {
        event.preventDefault()

        const res = await fetch(
            '/add',
            {
                body: JSON.stringify({
                    name: event.target.name.value,
                    surname: event.target.surname.value,
                }),
                headers: {
                    'Content-Type': 'application/json'
                },
                method: 'POST'
            }
        )

        const result = await res.json()
    }

    function Last10People() {
        const { data } = useSWR('/all', fetcherJson, { refreshInterval: 10000 })
        console.dir(data)

        return (<div>
            {data && data.map(todo => (
                <div key={todo.id}>id: {todo.id} name: {todo.name +" surname: "+ todo.surname}</div>
            ))}
        </div>);
    }

    return (
        <div>
            <h1>Foo</h1>
            <p>
                This is <code>pages/foo/index.tsx</code>.
            </p>
            <p>
                Check out <Link href="/foo/bar">bar</Link>.
            </p>

            <form onSubmit={registerPeople}>
                <label htmlFor="name">Name</label>
                <input id="name" name="name" type="text" autoComplete="name" value="some name" required /> <br/>
                <label htmlFor="surname">Surname</label>
                <input id="surname" name="surname" type="text" autoComplete="surname" value="some surname" required /> <br/>
                <button type="submit">Register</button>
            </form>

            <h2 className={styles.error}>Last 10 name</h2>
            <Last10People/>
        </div>
    );
}

export default Foo;
