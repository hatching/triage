// Copyright (C) 2019-2020 Hatching B.V.
// All rights reserved.

package types

import "time"

// static.json
type (
	StaticReport struct {
		Version string `json:"version"`

		Sample   ReportSample   `json:"sample"`
		Task     ReportTask     `json:"task"`
		Analysis ReportAnalysis `json:"analysis"`

		Signatures  []*Signature  `json:"signatures,omitempty"`
		Files       []*FileReport `json:"files"`
		UnpackCount int           `json:"unpack_count"`
		ErrorCount  int           `json:"error_count"`
		CompatKind  string        `json:"kind,omitempty"`

		Errors    []ReportedFailure `json:"errors,omitempty"`
		Extracted []*Extract        `json:"extracted,omitempty"`
	}
	ReportSample struct {
		ID        string `json:"sample"`
		Kind      string `json:"kind,omitempty"`
		Size      uint64 `json:"size,omitempty"`
		Target    string `json:"target,omitempty"`
		Submitted string `json:"submitted,omitempty"`
	}
	ReportTask struct {
		ID     string `json:"task"`
		Target string `json:"target,omitempty"`
	}
	ReportAnalysis struct {
		Reported string   `json:"reported,omitempty"`
		Score    int      `json:"score,omitempty"`
		Tags     []string `json:"tags,omitempty"`
	}
	FileReport struct {
		Name    string `json:"filename"`
		RelPath string `json:"relpath,omitempty"`
		Size    uint64 `json:"filesize"`
		Hashes
		Extensions []string `json:"exts"`
		Tags       []string `json:"tags"`
		Filetype   string   `json:"filetype,omitempty"`
		Mime       string   `json:"mime,omitempty"`
		Depth      int      `json:"depth"`
		Error      string   `json:"error,omitempty"`
		Kind       string   `json:"kind"`
		Selected   bool     `json:"selected"`
		RunAs      string   `json:"runas,omitempty"`
		Password   string   `json:"password,omitempty"`
	}
	Hashes struct {
		MD5    string `json:"md5,omitempty"`
		SHA1   string `json:"sha1,omitempty"`
		SHA256 string `json:"sha256,omitempty"`
		SHA512 string `json:"sha512,omitempty"`
		SSDeep string `json:"ssdeep,omitempty"`
	}
)

// triage_report.json
type (
	TriageReport struct {
		Version    string                 `json:"version"`
		Sample     TargetDesc             `json:"sample"`
		Task       TargetDesc             `json:"task"`
		Errors     []ReportedFailure      `json:"errors,omitempty"`
		Analysis   ReportAnalysisInfo     `json:"analysis,omitempty"`
		Processes  []Process              `json:"processes,omitempty"`
		Signatures []Signature            `json:"signatures"`
		Network    NetworkReport          `json:"network"`
		Debug      map[string]interface{} `json:"debug,omitempty"`
		Dumped     []Dump                 `json:"dumped,omitempty"`
		Extracted  []Extract              `json:"extracted,omitempty"`
	}
	TargetDesc struct {
		ID              string   `json:"id,omitempty"`
		CompatScore     int      `json:"score,omitempty"`
		Submitted       string   `json:"submitted,omitempty"`
		CompatCompleted string   `json:"completed,omitempty"`
		Target          string   `json:"target,omitempty"`
		Pick            string   `json:"pick,omitempty"`
		Type            string   `json:"type,omitempty"`
		Size            int64    `json:"size,omitempty"`
		MD5             string   `json:"md5,omitempty"`
		SHA1            string   `json:"sha1,omitempty"`
		SHA256          string   `json:"sha256,omitempty"`
		SHA512          string   `json:"sha512,omitempty"`
		SSDeep          string   `json:"ssdeep,omitempty"`
		Filetype        string   `json:"filetype,omitempty"`
		StaticTags      []string `json:"static_tags,omitempty"`
		UserTags        []string `json:"user_tags,omitempty"`
	}
	ReportedFailure struct {
		Task    string `json:"task,omitempty"`
		Backend string `json:"backend,omitempty"`
		Reason  string `json:"reason"`
	}
	ReportAnalysisInfo struct {
		Score          int      `json:"score,omitempty"`
		Tags           []string `json:"tags"`
		TTP            []string `json:"ttp,omitempty"`
		Features       []string `json:"features,omitempty"`
		Submitted      string   `json:"submitted,omitempty"`
		Reported       string   `json:"reported,omitempty"`
		MaxTimeNetwork int64    `json:"max_time_network,omitempty"`
		MaxTimeKernel  uint32   `json:"max_time_kernel,omitempty"`
		Backend        string   `json:"backend,omitempty"`
		Resource       string   `json:"resource,omitempty"`
		ResourceTags   []string `json:"resource_tags,omitempty"`
		Platform       string   `json:"platform,omitempty"`
	}
	Process struct {
		ProcID       int32       `json:"procid,omitempty"`
		ParentProcID int32       `json:"procid_parent,omitempty"`
		PID          uint64      `json:"pid"`
		PPID         uint64      `json:"ppid"`
		Cmd          interface{} `json:"cmd"`
		Image        string      `json:"image,omitempty"`
		Orig         bool        `json:"orig"`
		System       bool        `json:"-"`
		Started      uint32      `json:"started"`
		Terminated   uint32      `json:"terminated,omitempty"`
	}
	Signature struct {
		Label       string      `json:"label,omitempty"`
		Name        string      `json:"name"`
		Score       int         `json:"score,omitempty"`
		TTP         []string    `json:"ttp,omitempty"`
		Tags        []string    `json:"tags,omitempty"`
		Indicators  []Indicator `json:"indicators,omitempty"`
		YaraRule    string      `json:"yara_rule,omitempty"`
		Description string      `json:"desc,omitempty"`
		URL         string      `json:"url,omitempty"`
	}
	NetworkReport struct {
		Flows    []NetworkFlow        `json:"flows,omitempty"`
		Requests []NetworkRequest     `json:"requests,omitempty"`
		IPs      map[string]NetworkIP `json:"ips,omitempty"`
	}
	Dump struct {
		At     uint32 `json:"at"`
		PID    uint64 `json:"pid,omitempty"`
		ProcID int32  `json:"procid,omitempty"`
		Path   string `json:"path,omitempty"`
		Name   string `json:"name,omitempty"`
		Kind   string `json:"kind,omitempty"`
		Origin string `json:"origin,omitempty"`
		Addr   uint64 `json:"addr,omitempty"`
		Length uint64 `json:"length,omitempty"`
		MD5    string `json:"md5,omitempty"`
		SHA1   string `json:"sha1,omitempty"`
		SHA256 string `json:"sha256,omitempty"`
		SHA512 string `json:"sha512,omitempty"`
		SSDeep string `json:"ssdeep,omitempty"`
		Size   int    `json:"size,omitempty"`
		NSRL   string `json:"nsrl,omitempty"`
	}
	Extract struct {
		DumpedFile  string       `json:"dumped_file,omitempty"`
		Resource    string       `json:"resource,omitempty"`
		Config      *Config      `json:"config,omitempty"`
		Path        string       `json:"path,omitempty"`
		RansomNote  *Ransom      `json:"ransom_note,omitempty"`
		Dropper     *Dropper     `json:"dropper,omitempty"`
		Credentials *Credentials `json:"credentials,omitempty"`
	}
	Indicator struct {
		IOC          string `json:"ioc,omitempty"`
		Description  string `json:"description,omitempty"`
		At           uint32 `json:"at,omitempty"`
		SourcePID    uint64 `json:"pid,omitempty"`
		SourceProcID int32  `json:"procid,omitempty"`
		TargetPID    uint64 `json:"pid_target,omitempty"`
		TargetProcID int32  `json:"procid_target,omitempty"`
		Flow         int    `json:"flow,omitempty"`
		DumpFile     string `json:"dump_file,omitempty"`
		Resource     string `json:"resource,omitempty"`
		YaraRule     string `json:"yara_rule,omitempty"`
	}
	NetworkFlow struct {
		ID        int      `json:"id,omitempty"`
		Source    string   `json:"src,omitempty"`
		Dest      string   `json:"dst,omitempty"`
		Proto     string   `json:"proto,omitempty"`
		PID       uint64   `json:"pid,omitempty"`
		ProcID    int32    `json:"procid,omitempty"`
		FirstSeen int64    `json:"first_seen,omitempty"`
		LastSeen  int64    `json:"last_seen,omitempty"`
		RxBytes   uint64   `json:"rx_bytes,omitempty"`
		RxPackets uint64   `json:"rx_packets,omitempty"`
		TxBytes   uint64   `json:"tx_bytes,omitempty"`
		TxPackets uint64   `json:"tx_packets,omitempty"`
		Protocols []string `json:"protocols,omitempty"`
		Domain    string   `json:"domain,omitempty"`
		JA3       string   `json:"tls_ja3,omitempty"`
		JA3S      string   `json:"tls_ja3s,omitempty"`
		SNI       string   `json:"tls_sni,omitempty"`
	}
	NetworkIP struct {
		CIDR string `json:"cidr,omitempty"`
		CC   string `json:"cc,omitempty"`
		ASN  string `json:"asn,omitempty"`
	}
	NetworkRequest struct {
		Flow       int                    `json:"flow,omitempty"`
		Index      int                    `json:"index,omitempty"`
		At         uint32                 `json:"at,omitempty"`
		DomainReq  *NetworkDomainRequest  `json:"dns_request,omitempty"`
		DomainResp *NetworkDomainResponse `json:"dns_response,omitempty"`
		WebReq     *NetworkWebRequest     `json:"http_request,omitempty"`
		WebResp    *NetworkWebResponse    `json:"http_response,omitempty"`
	}
	Config struct {
		Family       string        `json:"family,omitempty"`
		Tags         []string      `json:"tags,omitempty"`
		Rule         string        `json:"rule,omitempty"`
		C2           []string      `json:"c2,omitempty"`
		Version      string        `json:"version,omitempty"`
		Botnet       string        `json:"botnet,omitempty"`
		Campaign     string        `json:"campaign,omitempty"`
		Mutex        []string      `json:"mutex,omitempty"`
		Decoy        []string      `json:"decoy,omitempty"`
		Wallet       []string      `json:"wallet,omitempty"`
		DNS          []string      `json:"dns,omitempty"`
		Keys         []Key         `json:"keys,omitempty"`
		Webinject    []string      `json:"webinject,omitempty"`
		CommandLines []string      `json:"command_lines,omitempty"`
		ListenAddr   string        `json:"listen_addr,omitempty"`
		ListenPort   int           `json:"listen_port,omitempty"`
		ListenFor    []string      `json:"listen_for,omitempty"`
		Shellcode    [][]byte      `json:"shellcode,omitempty"`
		ExtractedPE  []string      `json:"extracted_pe,omitempty"`
		Credentials  []Credentials `json:"credentials,omitempty"`
		Attributes   interface{}   `json:"attr,omitempty"`
	}
	Ransom struct {
		Family  string   `json:"family,omitempty"`
		Target  string   `json:"target,omitempty"`
		Emails  []string `json:"emails,omitempty"`
		Wallets []string `json:"wallets,omitempty"`
		URLs    []string `json:"urls,omitempty"`
		Contact []string `json:"contact,omitempty"`
		Note    string   `json:"note"`
	}
	Dropper struct {
		Family   string       `json:"family,omitempty"`
		Language string       `json:"language"`
		Source   string       `json:"source"`
		Deobf    string       `json:"deobfuscated"`
		URLs     []DropperURL `json:"urls"`
	}
	Credentials struct {
		Flow     int    `json:"flow,omitempty"`
		Protocol string `json:"protocol"`
		Host     string `json:"host,omitempty"`
		Port     int    `json:"port,omitempty"`
		User     string `json:"username"`
		Pass     string `json:"password"`
		EmailTo  string `json:"email_to,omitempty"`
	}
	NetworkDomainRequest struct {
		Domains   []string   `json:"domains,omitempty"`
		Questions []DNSEntry `json:"questions,omitempty"`
	}
	NetworkDomainResponse struct {
		Domains []string   `json:"domains,omitempty"`
		IP      []string   `json:"ip,omitempty"`
		Answers []DNSEntry `json:"answers,omitempty"`
	}
	NetworkWebRequest struct {
		Method  string   `json:"method,omitempty"`
		URL     string   `json:"url"`
		Request string   `json:"request"`
		Headers []string `json:"headers,omitempty"`
	}
	NetworkWebResponse struct {
		Status   string   `json:"status"`
		Response string   `json:"response"`
		Headers  []string `json:"headers,omitempty"`
	}
	Key struct {
		Kind  string      `json:"kind"`
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}
	DropperURL struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
	DNSEntry struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value,omitempty"`
	}
)

// overview.json
type (
	OverviewReport struct {
		Version    string              `json:"version"`
		Sample     OverviewSample      `json:"sample"`
		Tasks      []TaskSummary       `json:"tasks,omitempty"`
		Analysis   OverviewAnalysis    `json:"analysis"`
		Targets    []OverviewTarget    `json:"targets"`
		Errors     []ReportedFailure   `json:"errors,omitempty"`
		Signatures []Signature         `json:"signatures,omitempty"`
		Extracted  []OverviewExtracted `json:"extracted,omitempty"`
	}
	OverviewSample struct {
		TargetDesc
		Created   time.Time     `json:"created"`
		Completed time.Time     `json:"completed"`
		IOCs      *OverviewIOCs `json:"iocs,omitempty"`
	}
	TaskSummary struct {
		Sample     string   `json:"sample"`
		Kind       string   `json:"kind,omitempty"`
		Name       string   `json:"name,omitempty"`
		Status     string   `json:"status,omitempty"`
		TTP        []string `json:"ttp,omitempty"`
		Tags       []string `json:"tags,omitempty"`
		Score      int      `json:"score,omitempty"`
		Target     string   `json:"target,omitempty"`
		Backend    string   `json:"backend,omitempty"`
		Resource   string   `json:"resource,omitempty"`
		Platform   string   `json:"platform,omitempty"` // Deprecated
		TaskName   string   `json:"task_name,omitempty"`
		Failure    string   `json:"failure,omitempty"`
		Pick       string   `json:"pick,omitempty"`
		OS         string   `json:"os,omitempty"`
		Timeout    int      `json:"timeout,omitempty"`
		Signatures int      `json:"sigs,omitempty"`
	}
	OverviewAnalysis struct {
		Score  int      `json:"score"`
		Family []string `json:"family,omitempty"`
		Tags   []string `json:"tags,omitempty"`
	}
	OverviewTarget struct {
		Tasks []string `json:"tasks"`
		TargetDesc
		Tags       []string      `json:"tags,omitempty"`
		Family     []string      `json:"family,omitempty"`
		Signatures []Signature   `json:"signatures"`
		IOCs       *OverviewIOCs `json:"iocs,omitempty"`
	}
	OverviewExtracted struct {
		Tasks []string `json:"tasks"`
		*Extract
	}
	OverviewIOCs struct {
		URLs    []string `json:"urls,omitempty"`
		Domains []string `json:"domains,omitempty"`
		IPs     []string `json:"ips,omitempty"`
	}
)
