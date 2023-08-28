export interface PingResultResponse {
    results: { [key: string]: PingResult };
}

export interface PingResult {
    address: string;
    status: string;
    error: string;
    statistics: PingStatistics;
}

export interface PingStatistics {
    PacketsRecv: number;
    PacketsSent: number;
    PacketsRecvDuplicates: number;
    PacketLoss: number;
    Addr: string;
    Rtts: number[];
    MinRtt: number;
    MaxRtt: number;
    AvgRtt: number;
    StdDevRtt: number;
}