import {Button, Card, Col, DatePicker, Form, Input, Row, Select, Modal, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from "umi";
const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchRankList"],
    loadingDetail: state.loading.effects["global/fetchLogList"]
  })
)
class Rank extends React.Component {

  state = {
    isShowExpand: {},
    isModalVisible: false,
    searchPlayer: ''
  }

  renderName = (value) => {
    return <Link to={`/log?player=${value}`}>{value}</Link>
  }

  formRef = React.createRef();

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
        sorter: function (a, b) {
          return a.player.localeCompare(b.player)
        },
        render: this.renderName
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
        width: '10%',
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
        title: "技能占比",
        dataIndex: 'rate',
        key: 'rate',
        width: '10%',
        sorter: function (a, b) {
          return a.rate - b.rate
        },
        render: function (value) {
          return (value*100).toFixed(1) + "%"
        },
      },
      {
        title: "上榜次数",
        dataIndex: 'counts',
        key: 'counts',
        width: '10%',
        sorter: function (a, b) {
          return a.counts - b.counts
        },
        defaultSortOrder: "descend"
      },
      {
        title: "上榜时间点(最近30个)",
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
        width: 180,
        render: function (value) {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
        render: function (value, row) {
          let color = "grey"
          let typeName = ""
          if (row.player_type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.player_type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.player_type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <div><Tag className="custom-tag" color={color}>{typeName}</Tag><span>{value}</span></div>
        }
      },
      {
        title: "被玩家",
        dataIndex: 'target_player',
        key: 'target_player',
        render: function (value, row) {
          let color = "grey"
          let typeName = ""
          if (row.target_player_type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.target_player_type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.target_player_type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <span><Tag className="custom-tag" color={color}>{typeName}</Tag>{value}</span>
        }
      },
      {
        title: "伤害",
        dataIndex: 'damage',
        key: 'damage',
        width: 50
      },
      {
        title: "原始日志",
        dataIndex: 'origin_desc',
        key: 'origin_desc',
        render: function (value, row) {
          let results = []
          const parts = value.split(row.skill);
          results.push(parts[0])
          if (row.skill !== "普通攻击") {
            results.push(<span style={{color: "red", fontWeight: "bold"}} key={1}>{row.skill}</span>)
          }
          results.push(parts[1])
          return <div>{results}</div>;
        }
      },
    ];
  }

  handleUp = (player, action) => {
    const {isShowExpand} = this.state;
    isShowExpand[player] = action
    this.setState({
      isShowExpand: isShowExpand
    });
  }

  renderTimes = (value, row) => {
    const {isShowExpand} = this.state;
    let times = value.split(',')
    if (times.length > 9 && !isShowExpand[row.player]) {
      times = times.slice(0, 9)
      return (
        <div>
          <div>
            {
              times.map((v) => {
                return (<Tag color="green" onClick={() => this.searchRank({time: v, player: row.player})} key={v}>
                  {moment(v).format("YYYY-MM-DD HH:mm:ss")}</Tag>)
              })
            }<span onClick={() => this.handleUp(row.player, true)} style={{fontSize: '12px', color: 'orange'}}>展开</span>
          </div>
        </div>)
    } else {
      return (
        <div>
          <div>
            {
              times.map((v) => {
                return (<Tag color="green" onClick={() => this.searchRank({time: v, player: row.player})} key={v}>
                  {moment(v).format("YYYY-MM-DD HH:mm:ss")}</Tag>)
              })
            }{isShowExpand[row.player] &&
          <span onClick={() => this.handleUp(row.player, false)} style={{fontSize: '12px', color: 'orange'}}>收起</span>}
          </div>
        </div>)
    }
  }

  searchRank(record) {
    const {dispatch} = this.props
    this.setState({isModalVisible: true, searchPlayer: record.player})
    dispatch({
      type: 'global/fetchLogList',
      payload: {
        st: record.time,
        et: record.time,
        player: record.player
      },
    });
  }

  componentDidMount() {
    this.formRef.current.setFieldsValue({level: "4"})
    this.query()
  }

  query = (d, ds) => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    let st = fieldValue.time && fieldValue.time[0].format("YYYY-MM-DD HH:mm:ss")
    let et = fieldValue.time && fieldValue.time[1].format("YYYY-MM-DD HH:mm:ss")
    dispatch({
      type: 'global/fetchRankList',
      payload: {
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
        level: fieldValue.level ?? "4",
        name: fieldValue.name
      },
    });
  }

  onReset = () => {
    this.formRef.current.resetFields();
    this.formRef.current.setFieldsValue({level: "4"})
    this.query()
  };

  searchForm() {
    const dateFormat = 'YYYY-MM-DD HH:mm:ss';
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
        <Form.Item label="时间" name="time" style={{marginTop: "5px"}}>
          <RangePicker
            format={dateFormat}
            ranges={{
              今天: [moment().startOf('day'), moment().endOf('day')],
              昨天: [moment().subtract(1, 'day').startOf('day'), moment().subtract(1, 'day').endOf('day')],
              最近3天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近7天: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={(d, ds) => this.query(d, ds)}
          />
        </Form.Item>
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
            <Option value="3">至尊星耀</Option>
            <Option value="4">最强王者</Option>
            <Option value="5">荣耀王者</Option>
          </Select>
        </Form.Item>
        <Form.Item label="玩家" name="name" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
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
    const {rankList, loading, logList, loadingDetail} = this.props
    const {isModalVisible,searchPlayer} = this.state
    const listData = logList.list && logList.list.filter(v => {
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
              defaultPageSize: 15,
              hideOnSinglePage: true,
              showTotal: (total) => `共${total}条记录`,
            }}
            loading={loading}
          />
        </Card>
        <Modal
          title="日志详情"
          visible={isModalVisible}
          onCancel={()=>{this.setState({isModalVisible: false})}}
          width="70%"
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
              defaultPageSize: 10,
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