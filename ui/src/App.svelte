<script lang="ts">
    import {onMount} from 'svelte';
    import Button, {Icon, Label} from '@smui/button';
    import DataTable, {Body, Cell, Head, Row} from '@smui/data-table';
    import Paper, {Content, Title} from '@smui/paper';
    import Select, {Option} from '@smui/select';
    import moment from 'moment';
    import prettyBytes from 'pretty-bytes';
    import axios from 'axios';
    import {FlatToast, ToastContainer, toasts} from 'svelte-toasts';

    import type {VpnConnInfo} from './lib/strongswan.js';
    import type {PingResult} from './lib/ping';

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
    let connections: VpnConnInfo[] | null = null;
    let refreshTimer;
    let pinging = false;
    let pingResults: { [key: string]: PingResult } | null = null;

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
                connections = response.data.entries;
            })
            .catch(e => {
                connections = null;
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

<div class="main-body">
    <div>
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
                {#if connections}
                    {#each connections as connInfo}
                        <Row>
                            <Cell>
                                {connInfo['remote-id']}
                            </Cell>
                            <Cell>
                                {connInfo.IkeSaName}
                            </Cell>
                            <Cell>
                                {connInfo['remote-ts']}
                            </Cell>
                            <Cell>
                                {toDurationString(connInfo.established)}
                            </Cell>
                            <Cell>
                                {toDurationString(connInfo['install-time'])}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(connInfo['bytes-in'])}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(connInfo['packets-in'])}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(connInfo['bytes-out'])}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(connInfo['packets-out'])}
                            </Cell>
                        </Row>
                    {/each}
                {/if}
                </Body>
            </DataTable>
        </Content>
    </Paper>

    {#if pingResults}
        <Paper variant="unelevated">
            <Title>strongSwan connections</Title>
            <Content>
                <DataTable style="width: 100%;">
                    <Head>
                        <Row>
                            <Cell>IP address</Cell>
                            <Cell numeric>packets sent</Cell>
                            <Cell numeric>packets received</Cell>
                            <Cell numeric>packet loss (%)</Cell>
                            <Cell numeric>min (ms)</Cell>
                            <Cell numeric>avg (ms)</Cell>
                            <Cell numeric>max (ms)</Cell>
                            <Cell numeric>mdev (ms)</Cell>
                        </Row>
                    </Head>
                    <Body>
                    {#if pingResults.results}
                        {#each Object.values(pingResults.results) as entry}
                            <Row>
                                <Cell>
                                    {entry.address}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toLocaleInt(entry.statistics.PacketsSent)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toLocaleInt(entry.statistics.PacketsRecv)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toLocaleInt(entry.statistics.PacketLoss)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toPingDurationString(entry.statistics.MinRtt)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toPingDurationString(entry.statistics.AvgRtt)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toPingDurationString(entry.statistics.MaxRtt)}
                                    {/if}
                                </Cell>
                                <Cell numeric>
                                    {#if entry.statistics}
                                        {toPingDurationString(entry.statistics.StdDevRtt)}
                                    {/if}
                                </Cell>
                            </Row>
                        {/each}
                    {/if}
                    </Body>
                </DataTable>
            </Content>
        </Paper>
    {/if}
</div>

<ToastContainer placement="bottom-center" theme="dark" showProgress={true} let:data={data}>
    <FlatToast {data}/>
</ToastContainer>

<style>
    .main-body {
        margin: 1em;
    }
</style>
