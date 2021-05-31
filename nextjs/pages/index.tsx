import styles from '../styles/index.module.sass'
import Link from "next/link";
import useSWR from "swr";
import React from 'react';
import { Suspense } from 'react'



async function fetcher(url: string) {
  const resp = await fetch(url);
  return resp.text();
}

const fetcherJson = url => fetch(url).then(r => r.json())

function Index(): JSX.Element {
  const { data, error } = useSWR("/api", fetcher, { refreshInterval: 1000 });

  function Profile() {
    const { data } = useSWR('/all', fetcherJson, { refreshInterval: 0 })
    console.dir(data)

     return (<div>
              {data && data.map(todo => (
                <div key={todo.id}>{todo.Id} name: {todo.name +" "+ todo.surname}</div>
              ))}
            </div>);
          }

  return (
    <div className={styles.error}>
      <h1>Hello, world error!</h1>
      <p>
        This is <code>pages/index.tsx</code>.
      </p>
      <p>
        Check out <Link href="/foo">foo</Link>.
      </p>

      <h2>Memory allocation stats from Go server</h2>
      {error && (
        <p>
          Error fetching profile: <strong>{error}</strong>
        </p>
      )}
      {!error && !data && <p>Loading ...</p>}
      {!error && data && <pre>{data}</pre>}

        <Profile/>

      </div>
  );
}

export default Index;
