import React, { useState } from 'react'
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Paper from '@material-ui/core/Paper'
import Switch from '@material-ui/core/Switch'
import Link from "next/link";

const DarkTheme = () => {
    const [dark, setDark] = useState(false)

    const theme = createMuiTheme({
        palette: {
            type: dark ? 'dark' : 'light',
        },
    })

    return (
        <ThemeProvider theme={theme}>
            <Switch checked={dark} onChange={() => setDark(!dark)} />
            <Paper>
                <Typography variant='h1'>This is <code>pages/foo/darkTheme.tsx</code></Typography>

                <Typography variant='body2'>Check out <Link href="/">the homepage</Link>.</Typography>
            </Paper>
        </ThemeProvider>
    )
}

export default DarkTheme