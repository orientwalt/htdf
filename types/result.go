package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/orientwalt/htdf/codec"
	types "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	yaml "gopkg.in/yaml.v2"
)

// Result is the union of ResponseFormat and ResponseCheckTx.
type Result struct {
	// Code is the response code, is stored back on the chain.
	Code CodeType

	// Codespace is the string referring to the domain of an error
	Codespace CodespaceType

	// Data is any data returned from the app.
	// Data has to be length prefixed in order to separate
	// results from multiple msgs executions
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`

	// Log contains the txs log information. NOTE: nondeterministic.
	Log string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`

	// GasWanted is the maximum units of work we allow this tx to perform.
	GasWanted uint64

	// GasUsed is the amount of gas actually consumed. NOTE: unimplemented
	GasUsed uint64

	// Tags are used for transaction indexing and pubsub.
	// Tags Tags

	Events []types.Event `protobuf:"bytes,3,rep,name=events,proto3" json:"events"`
}

// TODO: In the future, more codes may be OK.
func (res Result) IsOK() bool {
	return res.Code.IsOK()
}

func (r Result) String() string {
	bz, _ := yaml.Marshal(r)
	return string(bz)
}

func (r Result) GetEvents() Events {
	events := make(Events, len(r.Events))
	for i, e := range r.Events {
		events[i] = Event(e)
	}

	return events
}

func (m *Result) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Result) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Result) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Log) > 0 {
		i -= len(m.Log)
		copy(dAtA[i:], m.Log)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Log)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Result) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Log)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *Result) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Result: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Result: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Log", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Log = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, types.Event{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure containing an indexed tx ABCI message log.
type ABCIMessageLog struct {
	MsgIndex int    `json:"msg_index"`
	Success  bool   `json:"success"`
	Log      string `json:"log"`
}

// String implements the fmt.Stringer interface for the ABCIMessageLogs type.
func (logs ABCIMessageLogs) String() (str string) {
	if logs != nil {
		raw, err := json.Marshal(logs)
		if err == nil {
			str = string(raw)
		}
	}

	return str
}

// ResultData represents the data returned in an sdk.Result
type ResultData struct {
	ContractAddress ethcmn.Address  `json:"contractAddress" gencodec:"required"`
	Bloom           ethtypes.Bloom  `json:"logsBloom" gencodec:"required"`
	Logs            []*ethtypes.Log `json:"logs" gencodec:"required"`
	Ret             []byte          `json:"ret" gencodec:"required"`
	TxHash          ethcmn.Hash     `json:"transactionHash" gencodec:"required"`
}

// String implements fmt.Stringer interface.
func (rd ResultData) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %0512x
	Logs: %v
	Ret: %s
	TxHash: %s
`, rd.ContractAddress.String(), fmt.Sprintf("0x%x", rd.Bloom.Big()), rd.Logs, fmt.Sprintf("0x%x", rd.Ret), rd.TxHash.String()))
}

func (rd ResultData) StringEx() string {
	return fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %0512x
	Logs: %v
	Ret: %s
	TxHash: %s
`, rd.ContractAddress.String(), fmt.Sprintf("0x%x", rd.Bloom.Big()), rd.Logs, fmt.Sprintf("0x%x", rd.Ret), rd.TxHash.String())
}

type ResultDataStr struct {
	ContractAddress string          `json:"contractAddress"`
	Bloom           string          `json:"logsBloom"`
	Logs            []*ethtypes.Log `json:"logs" gencodec:"required"`
	Ret             string          `json:"ret"`
	TxHash          string          `json:"transactionHash"`
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func NewResultDataStr(rd ResultData) ResultDataStr {
	var contractAddr string
	if !allZero(rd.ContractAddress.Bytes()) {
		contractAddr = rd.ContractAddress.String()
	}
	// var ret bool = false
	// if !allZero(rd.Ret) {
	// 	ret = true
	// }
	return ResultDataStr{
		ContractAddress: contractAddr,
		Bloom:           fmt.Sprintf("0x%x", rd.Bloom.Big()),
		Logs:            rd.Logs,
		Ret:             fmt.Sprintf("0x%x", rd.Ret),
		TxHash:          rd.TxHash.String(),
	}
}

func (rd ResultDataStr) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %s
	Logs: %s
	Ret: %s
	TxHash: %s
`, rd.ContractAddress, rd.Bloom, rd.Logs, rd.Ret, rd.TxHash))
}

type TxReceipt struct {
	Height    int64         `json:"height"`
	TxHash    string        `json:"txhash"`
	Results   ResultDataStr `json:"results,omitempty" gencodec:"required"`
	GasWanted int64         `json:"gas_wanted,omitempty"`
	GasUsed   int64         `json:"gas_used,omitempty"`
	Timestamp string        `json:"timestamp,omitempty"`
}

// EncodeResultData takes all of the necessary data from the EVM execution
// and returns the data as a byte slice encoded with amino
func EncodeResultData(data ResultData) ([]byte, error) {
	return codec.New().MarshalBinaryLengthPrefixed(data)
}

// DecodeResultData decodes an amino-encoded byte slice into ResultData
func DecodeResultData(in []byte) (ResultData, error) {
	var data ResultData
	err := codec.New().UnmarshalBinaryLengthPrefixed(in, &data)
	if err != nil {
		return ResultData{}, err
	}
	return data, nil
}

func NewResponseResultTxReceipt(res *ctypes.ResultTx, tx Tx, timestamp string) TxReceipt {
	if res == nil {
		return TxReceipt{}
	}

	data, err := DecodeResultData(res.TxResult.Data)
	if err != nil {
		return TxReceipt{}

	}
	return TxReceipt{
		TxHash:    res.Hash.String(),
		Height:    res.Height,
		Results:   NewResultDataStr(data),
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		Timestamp: timestamp,
	}
}

func (r TxReceipt) String() string {
	var sb strings.Builder
	sb.WriteString("Response:\n")

	if r.Height > 0 {
		sb.WriteString(fmt.Sprintf("  Height: %d\n", r.Height))
	}
	if r.TxHash != "" {
		sb.WriteString(fmt.Sprintf("  TxHash: %s\n", r.TxHash))
	}
	if r.Results.String() != "" {
		sb.WriteString(fmt.Sprintf("  Results: %s\n", r.Results.String()))
	}
	if r.GasWanted != 0 {
		sb.WriteString(fmt.Sprintf("  GasWanted: %d\n", r.GasWanted))
	}
	if r.GasUsed != 0 {
		sb.WriteString(fmt.Sprintf("  GasUsed: %d\n", r.GasUsed))
	}
	if r.Timestamp != "" {
		sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", r.Timestamp))
	}

	return strings.TrimSpace(sb.String())
}

// Empty returns true if the response is empty
func (r TxReceipt) Empty() bool {
	return r.TxHash == ""
}

// TxResponse defines a structure containing relevant tx data and metadata. The
// tags are stringified and the log is JSON decoded.
type TxResponse struct {
	Height    int64           `json:"height"`
	TxHash    string          `json:"txhash"`
	Codespace string          `json:"codespace,omitempty"`
	Code      uint32          `json:"code,omitempty"`
	Data      string          `json:"data,omitempty"`
	RawLog    string          `json:"raw_log,omitempty"`
	Logs      ABCIMessageLogs `json:"logs,omitempty"`
	Info      string          `json:"info,omitempty"`
	GasWanted int64           `json:"gas_wanted,omitempty"`
	GasUsed   int64           `json:"gas_used,omitempty"`
	Events    StringEvents    `json:"events,omitempty"`
	Tx        Tx              `json:"tx,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

// NewResponseResultTx returns a TxResponse given a ResultTx from tendermint
func NewResponseResultTx(res *ctypes.ResultTx, tx Tx, timestamp string) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	parsedLogs, _ := ParseABCILogs(res.TxResult.Log)

	return TxResponse{
		TxHash:    res.Hash.String(),
		Height:    res.Height,
		Codespace: res.TxResult.Codespace,
		Code:      res.TxResult.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.TxResult.Data)),
		RawLog:    res.TxResult.Log,
		Logs:      parsedLogs,
		Info:      res.TxResult.Info,
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		// Tags:      TagsToStringTags(res.TxResult.Tags),
		Tx:        tx,
		Timestamp: timestamp,
	}
}

// NewResponseFormatBroadcastTxCommit returns a TxResponse given a
// ResultBroadcastTxCommit from tendermint.
func NewResponseFormatBroadcastTxCommit(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	if !res.CheckTx.IsOK() {
		return newTxResponseCheckTx(res)
	}

	return newTxResponseDeliverTx(res)
}

func newTxResponseCheckTx(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.CheckTx.Log)

	return TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Codespace: res.CheckTx.Codespace,
		Code:      res.CheckTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.CheckTx.Data)),
		RawLog:    res.CheckTx.Log,
		Logs:      parsedLogs,
		Info:      res.CheckTx.Info,
		GasWanted: res.CheckTx.GasWanted,
		GasUsed:   res.CheckTx.GasUsed,
		// Tags:      TagsToStringTags(res.CheckTx.Tags),
	}
}

func newTxResponseDeliverTx(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.DeliverTx.Log)

	return TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Codespace: res.DeliverTx.Codespace,
		Code:      res.DeliverTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)),
		RawLog:    res.DeliverTx.Log,
		Logs:      parsedLogs,
		Info:      res.DeliverTx.Info,
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    StringifyEvents(res.DeliverTx.Events),
	}
}

// NewResponseFormatBroadcastTx returns a TxResponse given a ResultBroadcastTx from tendermint
func NewResponseFormatBroadcastTx(res *ctypes.ResultBroadcastTx) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	parsedLogs, _ := ParseABCILogs(res.Log)

	return TxResponse{
		Code:   res.Code,
		Data:   res.Data.String(),
		RawLog: res.Log,
		Logs:   parsedLogs,
		TxHash: res.Hash.String(),
	}
}

func (r TxResponse) String() string {
	var sb strings.Builder
	sb.WriteString("Response:\n")

	if r.Height > 0 {
		sb.WriteString(fmt.Sprintf("  Height: %d\n", r.Height))
	}
	if r.TxHash != "" {
		sb.WriteString(fmt.Sprintf("  TxHash: %s\n", r.TxHash))
	}
	if r.Code > 0 {
		sb.WriteString(fmt.Sprintf("  Code: %d\n", r.Code))
	}
	if r.Data != "" {
		sb.WriteString(fmt.Sprintf("  Data: %s\n", r.Data))
	}
	if r.RawLog != "" {
		sb.WriteString(fmt.Sprintf("  Raw Log: %s\n", r.RawLog))
	}
	if r.Logs != nil {
		sb.WriteString(fmt.Sprintf("  Logs: %s\n", r.Logs))
	}
	if r.Info != "" {
		sb.WriteString(fmt.Sprintf("  Info: %s\n", r.Info))
	}
	if r.GasWanted != 0 {
		sb.WriteString(fmt.Sprintf("  GasWanted: %d\n", r.GasWanted))
	}
	if r.GasUsed != 0 {
		sb.WriteString(fmt.Sprintf("  GasUsed: %d\n", r.GasUsed))
	}
	if len(r.Events) > 0 {
		sb.WriteString(fmt.Sprintf("  Events: \n%s\n", r.Events.String()))
	}

	if r.Codespace != "" {
		sb.WriteString(fmt.Sprintf("  Codespace: %s\n", r.Codespace))
	}
	if r.Timestamp != "" {
		sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", r.Timestamp))
	}

	return strings.TrimSpace(sb.String())
}

// Empty returns true if the response is empty
func (r TxResponse) Empty() bool {
	return r.TxHash == "" && r.Logs == nil
}

// SearchTxsResult defines a structure for querying txs pageable
type SearchTxsResult struct {
	TotalCount int          `json:"total_count"` // Count of all txs
	Count      int          `json:"count"`       // Count of txs in current page
	PageNumber int          `json:"page_number"` // Index of current page, start from 1
	PageTotal  int          `json:"page_total"`  // Count of total pages
	Limit      int          `json:"limit"`       // Max count txs per page
	Txs        []TxResponse `json:"txs"`         // List of txs in current page
}

func NewSearchTxsResult(totalCount, count, page, limit int, txs []TxResponse) SearchTxsResult {
	return SearchTxsResult{
		TotalCount: totalCount,
		Count:      count,
		PageNumber: page,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(limit))),
		Limit:      limit,
		Txs:        txs,
	}
}

// ParseABCILogs attempts to parse a stringified ABCI tx log into a slice of
// ABCIMessageLog types. It returns an error upon JSON decoding failure.
func ParseABCILogs(logs string) (res ABCIMessageLogs, err error) {
	err = json.Unmarshal([]byte(logs), &res)
	return res, err
}

// var _, _ types.UnpackInterfacesMessage = SearchTxsResult{}, TxResponse{}

// // UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// //
// // types.UnpackInterfaces needs to be called for each nested Tx because
// // there are generally interfaces to unpack in Tx's
// func (s SearchTxsResult) UnpackInterfaces(unpacker types.AnyUnpacker) error {
// 	for _, tx := range s.Txs {
// 		err := types.UnpackInterfaces(tx, unpacker)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// func (r TxResponse) UnpackInterfaces(unpacker types.AnyUnpacker) error {
// 	return types.UnpackInterfaces(r.Tx, unpacker)
// }
