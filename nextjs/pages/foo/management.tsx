import styles from '../../styles/index.module.sass'
import Link from "next/link";
import useSWR from "swr";
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import { DataGrid, GridRowsProp, GridColDef } from '@material-ui/data-grid';



const useStyles = makeStyles((theme) => ({
    root: {
        '& .MuiTextField-root': {
            margin: theme.spacing(1),
            width: '25ch',
        },
    },
}));

const fetcherJson = url => fetch(url).then(r => r.json())

function Management(): JSX.Element {

    const classes = useStyles();

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

        return !data?[]:data
    }

    return (
        <div>
            <h1>Foo</h1>
            <p>
                This is <code>pages/foo/index.tsx</code>.
            </p>
            <p>
                Check out <Link href="/foo/darkTheme" replace>Dark Theme</Link>.
            </p>

            <form onSubmit={registerPeople} className={classes.root} noValidate autoComplete="off">
                <TextField id="name" label="Name" defaultValue="some name" autoComplete="name" required /><p/>
                <TextField id="surname" label="Surname" defaultValue="some surname" autoComplete="surname" required /><p/>
                <Button variant="contained" color="primary" type="submit">
                    Register
                </Button>
            </form>

            <h2 className={styles.error}>Last 10 name</h2>
            <div style={{ height: 600, width: '100%' }}>
                <DataGrid rows={Last10People()}
                          columns={[
                                    { field: 'id', headerName: 'ID', width: 150 },
                                    { field: 'name', headerName: 'Name', width: 200 },
                                    { field: 'surname', headerName: 'Surname', width: 200 },
                                    { field: 'custom', headerName: 'Custom', width: 300 ,   valueGetter: (params) =>
                                            `${params.getValue(params.id, 'id') || ''} ${params.getValue(params.id, 'name') || ''}`,},
                                    ]}
                          hideFooterPagination={true}
                          rowHeight={32}/>
            </div>
        </div>
    );
}

export default Management;
