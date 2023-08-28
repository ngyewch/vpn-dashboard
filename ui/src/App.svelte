<script>
    import {onMount} from 'svelte';
    import Button, {Icon, Label} from '@smui/button';
    import DataTable, {Body, Cell, Head, Row} from '@smui/data-table';
    import Paper, {Content, Title} from '@smui/paper';
    import Select, {Option} from '@smui/select';
    import moment from 'moment';
    import prettyBytes from 'pretty-bytes';
    import axios from 'axios';
    import {FlatToast, ToastContainer, toasts} from 'svelte-toasts';

    const refreshOptions = [
        {
            value: -1,
            caption: "manual",
        }, {
            value: 5,
            caption: "5s",
        }, {
            value: 10,
            caption: "10s",
        }, {
            value: 15,
            caption: "15s",
        }, {
            value: 30,
            caption: "30s",
        }, {
            value: 60,
            caption: "1m",
        }
    ];

    let selectedRefreshOption = 5;
    let refreshing = false;
    let strongswanEntries = null;
    let refreshTimer;
    let pinging = false;
    let pingResults = null;

    $: {
        if (refreshTimer) {
            clearInterval(refreshTimer);
            refreshTimer = null;
        }
        if (selectedRefreshOption > 0) {
            refreshTimer = setInterval(() => {
                refresh();
            }, selectedRefreshOption * 1000);
        }
    }

    onMount(() => {
        refresh();
    });

    function ping() {
        pinging = true;
        pingResults = null;
        axios.get('/service/ping')
            .then(response => {
                const pingId = response.data.id;
                setTimeout(() => {
                    axios.get('/service/pingResult?id=' + pingId)
                        .then(response2 => {
                            pinging = false;
                            pingResults = response2.data;
                        })
                        .catch(e => {
                            pinging = false;
                            pingResults = null;
                            showError('Ping', e);
                        });
                }, 5000);
            })
            .catch(e => {
                pinging = false;
                pingResults = null;
                showError('Ping', e);
            });
    }

    function showError(title, e) {
        toasts.add({
            type: 'error',
            title: title,
            description: (e.response && e.response.data && e.response.data.error ? e.response.data.error : e),
        });
    }

    function refresh() {
        refreshing = true;
        const strongswanPromise = axios.get('/service/strongswan/connections')
            .then(response => {
                strongswanEntries = response.data.entries;
            })
            .catch(e => {
                strongswanEntries = null;
                showError('Refresh', e);
            });
        Promise.all([strongswanPromise])
            .then(responses => {
                refreshing = false;
            })
            .catch(e => {
                refreshing = false;
            });
    }

    function toDurationString(v) {
        if (isNaN(v)) {
            return null;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return null;
        }
        return moment.duration(n, 'seconds').humanize();
    }

    function toPingDurationString(v) {
        if (isNaN(v)) {
            return null;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return null;
        }
        return (n / 1000000).toLocaleString();
    }

    function toLocaleInt(v) {
        if (isNaN(v)) {
            return null;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return null;
        }
        return n.toLocaleString();
    }

    function toPrettyBytes(v) {
        if (isNaN(v)) {
            return null;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return null;
        }
        return prettyBytes(n);
    }
</script>

<style>
    .main-body {
        margin: 1em;
    }

    .align-right {
        text-align: right;
    }

    .table th.align-right {
        text-align: right;
    }
</style>

<div class="main-body">
    <div class="level">
        <Button variant="raised" on:click={(e) => ping()} disabled={pinging}>
            <Label>Ping</Label>
        </Button>
        <Button variant="raised" on:click={(e) => refresh()} disabled={refreshing}>
            <Icon class="material-icons">refresh</Icon>
            <Label>Refresh</Label>
        </Button>
        <Select bind:value={selectedRefreshOption}>
            {#each refreshOptions as refreshOption}
                <Option value={refreshOption.value}>{refreshOption.caption}</Option>
            {/each}
        </Select>
    </div>

    <Paper variant="unelevated">
        <Title>strongSwan connections</Title>
        <Content>
            <DataTable style="width: 100%;">
                <Head>
                    <Row>
                        <Cell>Remote ID</Cell>
                        <Cell>IKE SA name</Cell>
                        <Cell>Remote TS</Cell>
                        <Cell>Established</Cell>
                        <Cell>Installed</Cell>
                        <Cell numeric>Bytes in</Cell>
                        <Cell numeric>Packets in</Cell>
                        <Cell numeric>Bytes out</Cell>
                        <Cell numeric>Packets out</Cell>
                    </Row>
                </Head>
                <Body>
                {#if strongswanEntries}
                    {#each strongswanEntries as entry}
                        <Row>
                            <Cell>
                                {entry['remote-id']}
                            </Cell>
                            <Cell>
                                {entry['IkeSaName']}
                            </Cell>
                            <Cell>
                                {entry['remote-ts']}
                            </Cell>
                            <Cell>
                                {toDurationString(entry['established'])}
                            </Cell>
                            <Cell>
                                {toDurationString(entry['install-time'])}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(entry['bytes-in'])}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(entry['packets-in'])}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(entry['bytes-out'])}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(entry['packets-out'])}
                            </Cell>
                        </Row>
                    {/each}
                {/if}

                </Body>
            </DataTable>
        </Content>
    </Paper>

    {#if pingResults}
        <h3 class="title is-3">Ping results</h3>
        <table class="table is-fullwidth">
            <thead>
            <tr>
                <th>IP address</th>
                <th class="align-right">packets sent</th>
                <th class="align-right">packets received</th>
                <th class="align-right">packet loss (%)</th>
                <th class="align-right">min (ms)</th>
                <th class="align-right">avg (ms)</th>
                <th class="align-right">max (ms)</th>
                <th class="align-right">mdev (ms)</th>
            </tr>
            </thead>
            <tbody>
            {#if pingResults.results}
                {#each Object.values(pingResults.results) as entry}
                    <tr>
                        <td>
                            {entry.address}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toLocaleInt(entry.statistics.PacketsSent)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toLocaleInt(entry.statistics.PacketsRecv)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toLocaleInt(entry.statistics.PacketLoss)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toPingDurationString(entry.statistics.MinRtt)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toPingDurationString(entry.statistics.AvgRtt)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toPingDurationString(entry.statistics.MaxRtt)}
                            {/if}
                        </td>
                        <td class="align-right">
                            {#if entry.statistics}
                                {toPingDurationString(entry.statistics.StdDevRtt)}
                            {/if}
                        </td>
                    </tr>
                {/each}
            {/if}
            </tbody>
        </table>
    {/if}
</div>

<ToastContainer placement="bottom-center" theme="dark" showProgress={true} let:data={data}>
    <FlatToast {data}/>
</ToastContainer>
