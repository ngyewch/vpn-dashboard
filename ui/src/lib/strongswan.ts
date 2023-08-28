export interface VpnConnInfo {
    // IkeSa
    uniqueid: string;
    version: string;
    state: string;
    "local-host": string;
    "local-port": string;
    "local-id": string;
    "remote-host": string;
    "remote-port": string;
    "remote-id": string;
    "remote-xauth-id": string;
    "remote-eap-id": string;
    initiator: string;
    "initiator-spi": string;
    "responder-spi": string;
    "encr-alg": string;
    "encr-keysize": string;
    "integ-alg": string;
    "integ-keysize": string;
    "prf-alg": string;
    "dh-group": string;
    established: string;
    "rekey-time": string;
    "reauth-time": string;
    "remote-vips": string[];
    "child-sas": { [key: string] : Child_sas };
    "tasks-active": string[];
    "tasks-queued": string[];

    // Child_sas
    reqid: string;
    mode: string;
    protocol: string;
    encap: string;
    "spi-in": string;
    "spi-out": string;
    "cpi-in": string;
    "cpi-out": string;
    esn: string;
    "bytes-in": string;
    "packets-in": string;
    "use-in": string;
    "bytes-out": string;
    "packets-out": string;
    "use-out": string;
    "life-time": string;
    "install-time": string;
    "local-ts": string[];
    "remote-ts": string[];

    IkeSaName: string;
    ChildSaName: string;
}

export interface Child_sas {
    reqid: string;
    state: string;
    mode: string;
    protocol: string;
    encap: string;
    "spi-in": string;
    "spi-out": string;
    "cpi-in": string;
    "cpi-out": string;
    "encr-alg": string;
    "encr-keysize": string;
    "integ-alg": string;
    "integ-keysize": string;
    "prf-alg": string;
    "dh-group": string;
    esn: string;
    "bytes-in": string;
    "packets-in": string;
    "use-in": string;
    "bytes-out": string;
    "packets-out": string;
    "use-out": string;
    "rekey-time": string;
    "life-time": string;
    "install-time": string;
    "local-ts": string[];
    "remote-ts": string[];
}