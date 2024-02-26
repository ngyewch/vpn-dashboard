<script lang="ts">
    import {onMount} from 'svelte';
    import Button, {Icon, Label} from '@smui/button';
    import DataTable, {Body, Cell, Head, Row} from '@smui/data-table';
    import Paper, {Content, Title} from '@smui/paper';
    import Select, {Option} from '@smui/select';
    import moment from 'moment';
    import prettyBytes from 'pretty-bytes';
    import ky from 'ky';
    import {FlatToast, ToastContainer, toasts} from 'svelte-toasts';

    import type {GetConnectionsResponse, VpnConnInfo} from './lib/strongswan.js';

    interface RefreshOption {
        value: number;
        caption: string;
    }

    interface Entry {
        remoteId: string;
        ikeSaName: string;
        remoteTs: string[];
        established: number;
        bytesIn: number;
        packetsIn: number;
        bytesOut: number;
        packetsOut: number;
    }

    const refreshOptions: RefreshOption[] = [
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
    let entries: Entry[] | null = null;
    let refreshTimer: number | undefined;

    $: {
        if (refreshTimer) {
            clearInterval(refreshTimer);
            refreshTimer = undefined;
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

    function doParseInt(s: string | null | undefined): number {
        if ((s === null) || (s === undefined)) {
            return 0;
        }
        const v = parseInt(s);
        if (isNaN(v)) {
            return 0;
        }
        return v;
    }

    function refresh() {
        refreshing = true;
        ky.get('/service/strongswan/connections').json<GetConnectionsResponse>()
            .then(response => {
                refreshing = false;
                connections = response.entries;
                entries = connections.map(connection => {
                    let bytesIn = 0;
                    let packetsIn = 0;
                    let bytesOut = 0;
                    let packetsOut = 0;
                    if (connection["child-sas"]) {
                        for (const [key, childSas] of Object.entries(connection["child-sas"])) {
                            bytesIn += doParseInt(childSas["bytes-in"]);
                            packetsIn += doParseInt(childSas["packets-in"]);
                            bytesOut += doParseInt(childSas["bytes-out"]);
                            packetsOut += doParseInt(childSas["packets-out"]);
                        }
                    } else {
                        bytesIn = doParseInt(connection["bytes-in"]);
                        packetsIn = doParseInt(connection["packets-in"]);
                        bytesOut = doParseInt(connection["bytes-out"]);
                        packetsOut = doParseInt(connection["packets-out"]);
                    }
                    const entry: Entry = {
                        remoteId: connection['remote-id'],
                        ikeSaName: connection.IkeSaName,
                        remoteTs: connection['remote-ts'],
                        established: doParseInt(connection.established),
                        bytesIn: bytesIn,
                        packetsIn: packetsIn,
                        bytesOut: bytesOut,
                        packetsOut: packetsOut,
                    };
                    return entry;
                });
            })
            .catch(e => {
                refreshing = false;
                connections = null;
                entries = null;
                toasts.add({
                    type: 'error',
                    title: 'Refresh',
                    description: (e.response && e.response.data && e.response.data.error ? e.response.data.error : e),
                });
            });
    }

    function toDurationString(v: any, defaultValue: string): string {
        if (isNaN(v)) {
            return defaultValue;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return defaultValue;
        }
        return moment.duration(n, 'seconds').humanize();
    }

    function toLocaleInt(v: any, defaultValue: string): string {
        if (isNaN(v)) {
            return defaultValue;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return defaultValue;
        }
        return n.toLocaleString();
    }

    function toPrettyBytes(v: any, defaultValue: string): string {
        if (isNaN(v)) {
            return defaultValue;
        }
        const n = parseInt(v);
        if (isNaN(n)) {
            return defaultValue;
        }
        return prettyBytes(n);
    }
</script>

<div class="main-body">
    <div>
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
                        <Cell numeric>Bytes in</Cell>
                        <Cell numeric>Packets in</Cell>
                        <Cell numeric>Bytes out</Cell>
                        <Cell numeric>Packets out</Cell>
                    </Row>
                </Head>
                <Body>
                {#if entries}
                    {#each entries as entry}
                        <Row>
                            <Cell>
                                {entry.remoteId}
                            </Cell>
                            <Cell>
                                {entry.ikeSaName}
                            </Cell>
                            <Cell>
                                {entry.remoteTs}
                            </Cell>
                            <Cell>
                                {toDurationString(entry.established, '')}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(entry.bytesIn, '')}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(entry.packetsIn, '')}
                            </Cell>
                            <Cell numeric>
                                {toPrettyBytes(entry.bytesIn, '')}
                            </Cell>
                            <Cell numeric>
                                {toLocaleInt(entry.packetsOut, '')}
                            </Cell>
                        </Row>
                    {/each}
                {/if}
                </Body>
            </DataTable>
        </Content>
    </Paper>
</div>

<ToastContainer placement="bottom-center" theme="dark" showProgress={true} let:data={data}>
    <FlatToast {data}/>
</ToastContainer>

<style>
    .main-body {
        margin: 1em;
    }
</style>
