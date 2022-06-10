package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"nodemon/pkg/k8s"
	"nodemon/pkg/rpc"
	"strconv"
)

type Handler struct {
	RenNodeLabel string
	Ns           string
	Client       *kubernetes.Clientset
	Port         string
	MaxBlockDiff uint64
}

type GetAddressResponse struct {
	Name    string `json:"pod_name"`
	Address string `json:"address"`
	Error   string `json:"error"`
}

type CheckSyncResponse struct {
	Head  uint64 `json:"current_head"`
	Error string `json:"error"`
}

type GenericRPCResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

func (s *Handler) GetAddressHandler(c *gin.Context) {
	pods, err := k8s.GetPodIps(c, s.RenNodeLabel, s.Ns, s.Client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get pod ips, error=%v", err.Error()),
		})
		return
	}
	address := []GetAddressResponse{}
	for _, pod := range pods {
		if pod.IP == "" {
			address = append(address, GetAddressResponse{Error: "pod missing ip", Name: pod.Name})
			continue
		}
		addr, err := rpc.GetAddress(pod.IP, s.Port)
		if err != nil {
			address = append(address, GetAddressResponse{Error: err.Error(), Name: pod.Name})
			continue
		}
		address = append(address, GetAddressResponse{Address: addr.Hex(), Name: pod.Name})
	}

	c.JSON(http.StatusOK, address)
}

func (s *Handler) CheckSync(c *gin.Context) {
	// expects ?base=<uint>
	base, err := strconv.ParseUint(c.Query("base"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid base block number, error=%v", err.Error()),
		})
		return
	}
	pods, err := k8s.GetPodIps(c, s.RenNodeLabel, s.Ns, s.Client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get pod ips, error=%v", err.Error()),
		})
		return
	}
	address := map[string]CheckSyncResponse{}
	for _, pod := range pods {
		if pod.IP == "" {
			address[pod.Name] = CheckSyncResponse{Error: "pod missing ip"}
			continue
		}
		latest, err := rpc.GetLatestBlock(c, pod.IP, s.Port)
		if err != nil {
			address[pod.Name] = CheckSyncResponse{Error: err.Error()}
			continue
		}
		diff := latest - base
		if base > latest {
			diff = base - latest
		}
		if diff > s.MaxBlockDiff {
			address[pod.Name] = CheckSyncResponse{Head: latest, Error: fmt.Sprintf("diff from base(%v) %v > %v", base, diff, s.MaxBlockDiff)}
		}
	}

	c.JSON(http.StatusOK, address)
}

func (s *Handler) GenericRpc(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to read request body, error=%v", err.Error()),
		})
		return
	}

	// expects ?start=<uint>&end=<uint>
	pods, err := k8s.GetPodIps(c, s.RenNodeLabel, s.Ns, s.Client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get pod ips, error=%v", err.Error()),
		})
		return
	}
	start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid start index for pod, error=%v", err.Error()),
		})
		return
	}
	end, err := strconv.ParseInt(c.DefaultQuery("end", fmt.Sprintf("%v", len(pods)-1)), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid end index for pod, error=%v", err.Error()),
		})
		return
	}
	if start < 0 || start > end || end >= int64(len(pods)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid start and end index for pod, given start=%v end=%v, expected to be between %v and %v", start, end, 0, len(pods)-1),
		})
		return
	}
	address := map[string]GenericRPCResponse{}
	for i := start; i <= end; i++ {
		pod := pods[i]
		if pod.IP == "" {
			address[pod.Name] = GenericRPCResponse{Error: "pod missing ip"}
			continue
		}
		resp, err := rpc.GenericRpcCall(pod.IP, s.Port, data)
		if err != nil {
			address[pod.Name] = GenericRPCResponse{Error: err.Error()}
			continue
		}
		address[pod.Name] = GenericRPCResponse{Response: string(resp)}
	}

	c.JSON(http.StatusOK, address)
}
