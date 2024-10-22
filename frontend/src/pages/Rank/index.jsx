import {Button, Card, DatePicker, Form, Input, Modal, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from "umi";
import {getTypeColor, playerPros} from "@/utils/utils";

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchRankList"],
    loadingDetail: state.loading.effects["global/fetchLogData"]
  })
)
class Rank extends React.Component {

  state = {
    isModalVisible: false,
    searchPlayer: ''
  }
  formRef = React.createRef();

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
        width: "18%",
        sorter: function (a, b) {
          return a.player.localeCompare(b.player)
        },
        render: this.renderName
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
        width: '8%',
        sorter: function (a, b) {
          return a.type - b.type
        },
        render: function (value) {
          if (value === 0) {
            return <Tag color="orange">其它</Tag>
          }
          if (value === 1) {
            return <Tag color="green">天族</Tag>
          }
          if (value === 2) {
            return <Tag color="blue">魔族</Tag>
          }
        }
      },
      {
        title: "职业",
        dataIndex: 'class',
        key: 'class',
        width: '8%',
        sorter: function (a, b) {
          return a.class - b.class
        },
        render: function (value) {
          return <img src={require("../../assets/" + playerPros[value].logo)} width={35}/>
        }
      },
      {
        title: "上榜次数",
        dataIndex: 'counts',
        key: 'counts',
        width: '8%',
        sorter: function (a, b) {
          return a.counts - b.counts
        },
        defaultSortOrder: "descend"
      },
      {
        title: "上榜时间点(最近12个)",
        dataIndex: 'times',
        key: 'times',
        width: '50%',
        sorter: function (a, b) {
          return a.count - b.count
        },
        render: this.renderTimes
      },
    ];
    this.columnsDetail = [
      {
        title: "时间",
        dataIndex: 'time',
        key: 'time',
        width: 145,
        render: function (value) {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "原始日志",
        dataIndex: 'raw_msg',
        key: 'raw_msg',
        render: (value, row) => {
          if (row.skill === "kill" || row.skill === "killed") {
            return <div style={{color: "deeppink"}}>{value}</div>;
          }
          if (!row.skill) {
            return <div>{value}</div>;
          }
          const [color,typeName] = getTypeColor(row.player_type)
          let results = []
          const parts = value.split(row.player);
          results.push(parts[0])
          results.push((<span><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag
            className="custom-tag">{playerPros[row.player_class].name}</Tag><span>{row.player}</span></span>))

          const parts2 = parts[1].split(row.skill);
          results.push(parts2[0]);
          if (row.skill !== "attack") {
            results.push(<span style={{color: "red", fontWeight: "bold"}} key={1}>{row.skill}</span>)
          }
          if (row.target !== "" && parts2[1]) {
            const parts3 = parts2[1].split(row.target);
            const [color,typeName] = getTypeColor(row.target_type)
            results.push(parts3[0]);
            results.push((<span><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag
              className="custom-tag">{playerPros[row.target_class].name}</Tag><span>{row.target}</span></span>))
            results.push(parts3[1]);
          }else {
            results.push(parts2[1])
          }
          return <div>{results}</div>;
        },
      },
    ];
  }


  renderName = (value) => {
    return <Link to={`/log?player=${value}`}>{value}</Link>
  }

  renderTimes = (value, row) => {
    let times = value.split(',')
    if (times.length > 12) {
      times = times.slice(0, 12);
    }
    return (
      <div>
        <div>
          {
            times.map((v) => {
              return (<Tag color="green" onClick={() => this.searchRank({time: v, player: row.player})} key={v}>
                {moment(v).format("YYYY-MM-DD HH:mm:ss")}</Tag>)
            })
          }
        </div>
      </div>)
  }

  searchRank(record) {
    const {dispatch} = this.props
    this.setState({isModalVisible: true, searchPlayer: record.player})
    dispatch({
      type: 'global/fetchLogData',
      payload: {
        st: record.time,
        et: record.time,
        player: record.player
      },
    });
  }

  componentDidMount() {
    this.formRef.current.setFieldsValue({level: "3"})
    this.query()
  }

  query = () => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    dispatch({
      type: 'global/fetchRankList',
      payload: {
        level: fieldValue.level ?? "3",
        name: fieldValue.name,
        class: fieldValue.class
      },
    });
  }

  onReset = () => {
    this.formRef.current.resetFields();
    this.formRef.current.setFieldsValue({level: "3"})
    this.query()
  };

  searchForm() {
    const onFinish = async () => {
      this.query()
    };
    return (
      <Form
        layout="inline"
        onFinish={onFinish}
        autoComplete="false"
        ref={this.formRef}
      >
        <Form.Item label="段位" name="level" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            style={{width: 150}}
            placeholder="请选择段位"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            <Option value="3">黄金</Option>
            <Option value="4">钻石</Option>
            <Option value="5">王者</Option>
          </Select>
        </Form.Item>
        <Form.Item label="玩家" name="name" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
        </Form.Item>
        <Form.Item label="职业" name="class" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            style={{width: 150}}
            placeholder="请选择职业"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            {playerPros.map((v, k) =>
              <Option value={k} key={k}>{v.name}</Option>
            )}
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" style={{marginTop: "5px"}}>
            搜索
          </Button>
          &nbsp;&nbsp;
          <Button type="primary" onClick={this.onReset} style={{marginTop: "5px"}}>
            重置
          </Button>
        </Form.Item>
      </Form>)
  }

  render() {
    const {rankList, loading, logData, loadingDetail} = this.props
    const {isModalVisible, searchPlayer} = this.state
    const listData = logData.list && logData.list.filter(v => {
      return v.player === searchPlayer
    })
    return (
      <PageContainer>
        <Card extra={this.searchForm()} >
          <Table
            bordered
            size="small"
            columns={this.columns}
            dataSource={rankList}
            rowKey={(record) => {
              return record.time + record.player
            }}
            pagination={{
              defaultPageSize: 20,
              hideOnSinglePage: true,
              showTotal: (total) => `共${total}条记录`,
            }}
            loading={loading}
          />
        </Card>
        <Modal
          title="日志详情"
          visible={isModalVisible}
          onCancel={() => {
            this.setState({isModalVisible: false})
          }}
          width="60%"
          footer={null}
        >
          <Table
            bordered
            size="small"
            columns={this.columnsDetail}
            dataSource={listData}
            rowKey={(record) => {
              return record.id
            }}
            loading={loadingDetail}
            pagination={{
              defaultPageSize: 20,
              hideOnSinglePage: true,
              showTotal: (total) => `共${total}条记录`,
            }}
          />
        </Modal>
      </PageContainer>
    );
  }
}

export default Rank;
