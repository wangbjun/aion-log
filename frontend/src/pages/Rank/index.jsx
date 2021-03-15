import {Button, Card, Col, DatePicker, Form, Input, Row, Select, Statistic, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchRankList"]
  })
)
class Rank extends React.Component {

  state = {
    isShowExpand: {}
  }

  formRef = React.createRef();

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
        width: '30%',
        sorter: function (a, b) {
          return a.player.localeCompare(b.player)
        },
        render: function (value, row) {
          let color = "grey"
          let typeName = ""
          if (row.type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <span><Tag className="custom-tag" color={color}>{typeName}</Tag>{value}</span>
        }
      },
      {
        title: "上榜次数",
        dataIndex: 'count',
        key: 'count',
        width: '20%',
        sorter: function (a, b) {
          return a.count - b.count
        },
        defaultSortOrder: "descend"
      },
      {
        title: "上榜时间点",
        dataIndex: 'times',
        key: 'times',
        width: '45%',
        sorter: function (a, b) {
          return a.count - b.count
        },
        render: this.renderTimes
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
    if (times.length > 10 && !isShowExpand[row.player]) {
      times = times.slice(0, 10)
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

  async searchRank(record) {
    const {dispatch} = this.props
    await dispatch({
      type: 'global/saveDefault',
      payload: {
        sTime: record.time,
        sPlayer: record.player
      },
    });
  }

  componentDidMount() {
    this.formRef.current.setFieldsValue({
      time: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')]
    })
    this.formRef.current.setFieldsValue({level: "3"})
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
        level: fieldValue.level ?? "3",
        name: fieldValue.name
      },
    });
    dispatch({
      type: 'global/fetchStat',
      payload: {
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
      },
    });
  }

  onReset = () => {
    this.formRef.current.resetFields();
    this.formRef.current.setFieldsValue({
      time: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')]
    })
    this.formRef.current.setFieldsValue({level: "3"})
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
            <Option value="2">倔强青铜</Option>
            <Option value="3">荣耀黄金</Option>
            <Option value="4">最强王者</Option>
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
    const {rankList, stat, loading} = this.props
    return (
      <PageContainer>
        <Card title="概况" extra={this.searchForm()}>
          <Row gutter={24}>
            <Col span={8}>
              <Statistic title="玩家总数" value={stat.total}/>
            </Col>
            <Col span={8}>
              <Statistic title="上榜玩家" value={rankList ? rankList.length : 0}/>
            </Col>
            <Col span={8}>
              <Statistic title="上榜玩家占比" value={(rankList ? rankList.length / stat.total * 100 : 0).toFixed(2)}
                         suffix={"%"}/>
            </Col>
          </Row>
        </Card>
        <p/>
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
      </PageContainer>
    );
  }
}

export default Rank;
