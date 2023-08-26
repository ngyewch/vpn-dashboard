<script>
    import {onMount} from 'svelte';
    import {Button, Select, Toast} from 'svelma';
    import moment from 'moment';
    import prettyBytes from 'pretty-bytes';
    import axios from 'axios';

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
    let wireguardEntries = null;
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
        Toast.create({
            type: 'is-danger',
            position: 'is-bottom',
            message: title + ': '
                + (e.response && e.response.data && e.response.data.error ? e.response.data.error : e),
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
        const wireguardPromise = axios.get('/service/wireguard/connections')
            .then(response => {
                wireguardEntries = response.data.entries;
            })
            .catch(e => {
                wireguardEntries = null;
                showError('Refresh', e);
            });
        Promise.all([strongswanPromise, wireguardPromise])
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
        <div class="level-left">
        </div>
        <div class="level-right">
            <div class="level-item">
                <Button type="is-primary" on:click={(e) => ping()} loading={pinging} disabled={pinging}>
                    Ping
                </Button>
            </div>
            <div class="level-item">
                <Button type="is-primary" on:click={(e) => refresh()} loading={refreshing} disabled={refreshing}
                        iconPack="fas" iconRight="sync">
                    Refresh
                </Button>
            </div>
            <div class="level-item">
                <Select bind:selected={selectedRefreshOption}>
                    {#each refreshOptions as refreshOption}
                        <option value={refreshOption.value}>{refreshOption.caption}</option>
                    {/each}
                </Select>
            </div>
        </div>
    </div>
    <h3 class="title is-3">strongSwan connections</h3>
    <table class="table is-fullwidth">
        <thead>
        <tr>
            <th>Remote ID</th>
            <th>IKE SA name</th>
            <th>Remote TS</th>
            <th>Established</th>
            <th>Installed</th>
            <th class="align-right">Bytes in</th>
            <th class="align-right">Packets in</th>
            <th class="align-right">Bytes out</th>
            <th class="align-right">Packets out</th>
        </tr>
        </thead>
        <tbody>
        {#if strongswanEntries}
            {#each strongswanEntries as entry}
                <tr>
                    <td>
                        {entry['remote-id']}
                    </td>
                    <td>
                        {entry['IkeSaName']}
                    </td>
                    <td>
                        {entry['remote-ts']}
                    </td>
                    <td>
                        {toDurationString(entry['established'])}
                    </td>
                    <td>
                        {toDurationString(entry['install-time'])}
                    </td>
                    <td class="align-right">
                        {toPrettyBytes(entry['bytes-in'])}
                    </td>
                    <td class="align-right">
                        {toLocaleInt(entry['packets-in'])}
                    </td>
                    <td class="align-right">
                        {toPrettyBytes(entry['bytes-out'])}
                    </td>
                    <td class="align-right">
                        {toLocaleInt(entry['packets-out'])}
                    </td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
    <h3 class="title is-3">Wireguard connections</h3>
    <table class="table is-fullwidth">
        <thead>
        <tr>
            <th>Name</th>
            <th>Allowed IPs</th>
            <th>Endpoint</th>
            <th>Latest handshake</th>
            <th>Transfer</th>
        </tr>
        </thead>
        <tbody>
        {#if wireguardEntries}
            {#each wireguardEntries as entry}
                <tr>
                    <td>
                        {entry['name']}
                    </td>
                    <td>
                        {entry['allowedIps']}
                    </td>
                    <td>
                        {entry['endpoint']}
                    </td>
                    <td>
                        {entry['latestHandshake']}
                    </td>
                    <td>
                        {entry['transfer']}
                    </td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
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
