# Scaling ElastiCache for Memcached clusters

- from: https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/Scaling.html

- 어플리케이션의 데터용량을 고정하는 것은 드물다. 
- 비즈니스가 증가에 따라 필요에 따라 증감이 발생한다. 
- 만약 캐시를 자체 관리한다면, 충분한 하드웨어를 준비해야하고, 피크타임을 준비할 필요가 있다. 이는 매우 비싼 값이다. 
- Amazon ElastiCache 를 이용하여 현재 요구에 따라 활장할 수 있다. ElastiCache 는 필요에 맞게 스케일링 할 수있다. 

- 다음과 같은 클러스터 스케일링에 대한 사항이다. 
  - Scaling out: 클러스터를 노드에 추가 (https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/Clusters.AddNode.html)
  - Scaling in: 클러스터로 노드를 줄인다. (https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/Clusters.DeleteNode.html)
  - Changing node types: 수직으로 확장한다. (https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/Scaling.html#Scaling.Memcached.Vertically)

- Memcached 클러스터들은 1개에서 40개의 노드로 구성된다. 
- Memcached 클러스터를 스케일 아웃하는 것은 클러스터에 노드를 축하거나 제거하는 것으로 쉽게 할 수 있다. 
- 만약 40노드 이상이 필요하다면 AWS 리전에서 총 300 노드 이상을 원한다면 AWS에 요청하여 제한용량을 확장 요청하자. 

- 클러스터 용량을 증가하고자 한다면, 메모리 위주의 노드 타입으로 스케일 업 할수 있다. 
- Memcached 엔진은 노드를 데이터에 저장하지 않기 때문에 

## 수평확장 

- 수평 확장은 1개 ~ 40개 노드로 확장할 수 있다. 
- 쉽게 노드를 추가/삭제 하는 것으로 스케일 할 수 있다. 
- Memcached 클러스터에서 노드의 수를 변경하면, 올바른 노드에 매핑될 수 있도록 하기 위해서 키 스페이스를 다시 매핑해야한다. 
- auto discovery 를 사용하면 노드를 추가/삭제 하더라도 엔드포인트의 변경이 필요하지 않다. 
- 만약 auto discovery 를 사용하지 않는다면, 매번 노드 추가/삭제시 엔드포인트를 변경해야한다. 

## 수직확장 

- 노드 인스턴스 타입을 변경하여 메모리 확장등을 수행할 수 있다. 