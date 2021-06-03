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
      <h1 className={styles.error}>Hello, world!</h1>
      <p>
        This is <code>pages/index.tsx</code>.
      </p>
      <p>
        Check out <Link href="/foo">foo</Link>.
      </p>

      <h2 className={styles.error}>Last 10 name</h2>
      <Profile/>

      <h2 className={styles.error}>Memory allocation stats from Go server</h2>
      {error && (
        <p>
          Error fetching profile: <strong>{error}</strong>
        </p>
      )}
      {!error && !data && <p>Loading ...</p>}
      {!error && data && <pre>{data}</pre>}
    </div>
  );
}

export default Index;
