package utils

import (
	"fmt"
	"sync"
)

// MockAWSService provides a mock implementation of AWS services for testing
type MockAWSService struct {
	resources map[string]map[string]interface{}
	mutex     sync.RWMutex
}

// NewMockAWSService creates a new MockAWSService
func NewMockAWSService() *MockAWSService {
	return &MockAWSService{
		resources: make(map[string]map[string]interface{}),
	}
}

// AddResource adds a resource to the mock AWS service
func (m *MockAWSService) AddResource(resourceType, resourceID string, attributes map[string]interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.resources[resourceType]; !ok {
		m.resources[resourceType] = make(map[string]interface{})
	}
	m.resources[resourceType][resourceID] = attributes
}

// GetResource gets a resource from the mock AWS service
func (m *MockAWSService) GetResource(resourceType, resourceID string) (map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if resources, ok := m.resources[resourceType]; ok {
		if resource, ok := resources[resourceID]; ok {
			if attrs, ok := resource.(map[string]interface{}); ok {
				return attrs, nil
			}
		}
	}
	return nil, fmt.Errorf("resource %s of type %s not found", resourceID, resourceType)
}

// ListResources lists all resources of a given type
func (m *MockAWSService) ListResources(resourceType string) ([]map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if resources, ok := m.resources[resourceType]; ok {
		result := make([]map[string]interface{}, 0, len(resources))
		for _, resource := range resources {
			if attrs, ok := resource.(map[string]interface{}); ok {
				result = append(result, attrs)
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("no resources of type %s found", resourceType)
}

// DeleteResource deletes a resource from the mock AWS service
func (m *MockAWSService) DeleteResource(resourceType, resourceID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if resources, ok := m.resources[resourceType]; ok {
		if _, ok := resources[resourceID]; ok {
			delete(resources, resourceID)
			return nil
		}
	}
	return fmt.Errorf("resource %s of type %s not found", resourceID, resourceType)
}

// ClearResources clears all resources from the mock AWS service
func (m *MockAWSService) ClearResources() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.resources = make(map[string]map[string]interface{})
}

// MockVPC represents a mock VPC resource
type MockVPC struct {
	ID                string
	CIDRBlock         string
	EnableDNSSupport  bool
	EnableDNSHostnames bool
	Tags              map[string]string
}

// MockSubnet represents a mock Subnet resource
type MockSubnet struct {
	ID                    string
	VPCID                 string
	CIDRBlock             string
	AvailabilityZone      string
	MapPublicIPOnLaunch   bool
	Tags                  map[string]string
}

// MockInternetGateway represents a mock Internet Gateway resource
type MockInternetGateway struct {
	ID    string
	VPCID string
	Tags  map[string]string
}

// MockNATGateway represents a mock NAT Gateway resource
type MockNATGateway struct {
	ID            string
	SubnetID      string
	AllocationID  string
	Tags          map[string]string
}

// MockRouteTable represents a mock Route Table resource
type MockRouteTable struct {
	ID     string
	VPCID  string
	Routes []MockRoute
	Tags   map[string]string
}

// MockRoute represents a mock Route resource
type MockRoute struct {
	DestinationCIDRBlock string
	GatewayID            string
	NATGatewayID         string
}

// MockSecurityGroup represents a mock Security Group resource
type MockSecurityGroup struct {
	ID          string
	VPCID       string
	Name        string
	Description string
	IngressRules []MockSecurityGroupRule
	EgressRules  []MockSecurityGroupRule
	Tags         map[string]string
}

// MockSecurityGroupRule represents a mock Security Group Rule
type MockSecurityGroupRule struct {
	Protocol    string
	FromPort    int
	ToPort      int
	CIDRBlocks  []string
	Description string
}

// MockRDSInstance represents a mock RDS Instance resource
type MockRDSInstance struct {
	ID                 string
	DBName             string
	Engine             string
	EngineVersion      string
	InstanceClass      string
	AllocatedStorage   int
	StorageType        string
	SubnetGroupName    string
	SecurityGroupIDs   []string
	MultiAZ            bool
	PubliclyAccessible bool
	Tags               map[string]string
}

// MockElastiCacheCluster represents a mock ElastiCache Cluster resource
type MockElastiCacheCluster struct {
	ID                string
	ReplicationGroupID string
	Engine            string
	EngineVersion     string
	NodeType          string
	NumCacheNodes     int
	SubnetGroupName   string
	SecurityGroupIDs  []string
	Tags              map[string]string
}

// MockECSCluster represents a mock ECS Cluster resource
type MockECSCluster struct {
	ID    string
	Name  string
	Tags  map[string]string
}

// MockECSService represents a mock ECS Service resource
type MockECSService struct {
	ID                  string
	ClusterID           string
	TaskDefinitionARN   string
	DesiredCount        int
	LaunchType          string
	NetworkConfiguration MockECSNetworkConfiguration
	LoadBalancers       []MockECSLoadBalancer
	Tags                map[string]string
}

// MockECSNetworkConfiguration represents a mock ECS Network Configuration
type MockECSNetworkConfiguration struct {
	SubnetIDs         []string
	SecurityGroupIDs  []string
	AssignPublicIP    bool
}

// MockECSLoadBalancer represents a mock ECS Load Balancer
type MockECSLoadBalancer struct {
	TargetGroupARN string
	ContainerName  string
	ContainerPort  int
}

// MockALB represents a mock Application Load Balancer resource
type MockALB struct {
	ID            string
	Name          string
	SubnetIDs     []string
	SecurityGroupIDs []string
	Tags          map[string]string
}

// MockALBTargetGroup represents a mock ALB Target Group resource
type MockALBTargetGroup struct {
	ID        string
	Name      string
	Port      int
	Protocol  string
	VPCID     string
	TargetType string
	Tags      map[string]string
}

// MockALBListener represents a mock ALB Listener resource
type MockALBListener struct {
	ID             string
	LoadBalancerARN string
	Port           int
	Protocol       string
	DefaultAction  MockALBListenerAction
}

// MockALBListenerAction represents a mock ALB Listener Action
type MockALBListenerAction struct {
	Type           string
	TargetGroupARN string
}