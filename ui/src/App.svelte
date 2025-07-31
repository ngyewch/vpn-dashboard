<script lang="ts">
    import {onMount} from 'svelte';
    import prettyBytes from 'pretty-bytes';
    import ky from 'ky';
    import humanizeDuration from 'humanize-duration';
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

    const shortEnglishHumanizer = humanizeDuration.humanizer({
        language: "shortEn",
        languages: {
            shortEn: {
                y: () => "y",
                mo: () => "mo",
                w: () => "w",
                d: () => "d",
                h: () => "h",
                m: () => "m",
                s: () => "s",
                ms: () => "ms",
            },
        },
    });
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
        return shortEnglishHumanizer(v * 1000);
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
    <div class="toolbar">
        <button onclick={(e) => refresh()} disabled={refreshing}>Refresh</button>
        <select bind:value={selectedRefreshOption}>
            {#each refreshOptions as refreshOption}
                <option value={refreshOption.value}>{refreshOption.caption}</option>
            {/each}
        </select>
    </div>

    <h2>strongSwan connections</h2>
    <table>
        <thead>
        <tr>
            <th>Remote ID</th>
            <th>IKE SA name</th>
            <th>Remote TS</th>
            <th>Established</th>
            <th class="numeric">Bytes in</th>
            <th class="numeric">Packets in</th>
            <th class="numeric">Bytes out</th>
            <th class="numeric">Packets out</th>
        </tr>
        </thead>
        <tbody>
        {#if entries}
            {#each entries as entry}
                <tr>
                    <td>
                        {entry.remoteId}
                    </td>
                    <td>
                        {entry.ikeSaName}
                    </td>
                    <td>
                        {entry.remoteTs}
                    </td>
                    <td>
                        {toDurationString(entry.established, '')}
                    </td>
                    <td class="numeric">
                        {toPrettyBytes(entry.bytesIn, '')}
                    </td>
                    <td class="numeric">
                        {toLocaleInt(entry.packetsIn, '')}
                    </td>
                    <td class="numeric">
                        {toPrettyBytes(entry.bytesIn, '')}
                    </td>
                    <td class="numeric">
                        {toLocaleInt(entry.packetsOut, '')}
                    </td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
</div>

<ToastContainer placement="bottom-center" theme="light" showProgress={true} let:data={data}>
    <FlatToast {data}/>
</ToastContainer>

<style>
    .main-body {
        margin: 1em;
    }

    .toolbar {
        display: flex;
        gap: 0.5em;
        align-items: center;
        justify-content: center;
        margin-bottom: 2em;
    }

    .toolbar select {
        margin: 0;
    }

    th.numeric,
    td.numeric {
        text-align: right;
    }
</style>
