import {Button, Card, Col, DatePicker, Form, Input, Row, Select, Statistic, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from 'dva/router'

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchPlayerList"]
  })
)
class Player extends React.Component {
  formRef = React.createRef();

  renderName = (value) => {
    return <Link to={`/log?player=${value}`}>{value}</Link>
  }

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "玩家",
        dataIndex: 'name',
        key: 'name',
        defaultSortOrder: 'ascend',
        width: '20%',
        sorter: function (a, b) {
          return a.name.localeCompare(b.name)
        },
        render: this.renderName
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
        width: '20%',
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
        dataIndex: 'job',
        key: 'job',
        width: '20%',
        sorter: function (a, b) {
          return a.job.localeCompare(b.job)
        },
        render: function (value) {
          return "待完善"
        }
      },
      {
        title: "最后更新时间",
        dataIndex: 'time',
        key: 'time',
        width: '20%',
        sorter: function (a, b) {
          return moment(a.time).isAfter(moment(b.time))
        },
        render: function (value) {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "操作",
        dataIndex: 'option',
        key: 'option',
        width: '20%',
        render: this.renderOption
      },
    ];
  }

  renderOption = (value, row) => {
    return (<div>
      <a onClick={() => this.changeType(row, 1)}><Tag color="green">设为天族</Tag></a>
      <a onClick={() => this.changeType(row, 2)}><Tag color="blue">设为魔族</Tag></a>
      <a onClick={() => this.changeType(row, 0)}><Tag color="orange">设为其它</Tag></a>
    </div>)
  }

  async changeType(row, type) {
    const {dispatch} = this.props
    await dispatch({
      type: 'global/changePlayerType',
      payload: {
        id: row.id,
        type: type
      }
    });
    this.query()
  }

  componentDidMount() {
    this.query()
  }

  query = () => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    dispatch({
      type: 'global/fetchPlayerList',
      payload: {
        name: fieldValue.name,
        type: fieldValue.type
      }
    });
  }

  onReset = () => {
    this.formRef.current.resetFields();
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
        <Form.Item label="玩家" name="name" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
        </Form.Item>
        <Form.Item label="种族" name="type" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            style={{width: 150}}
            placeholder="请选择种族"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            <Option value="1">天族</Option>
            <Option value="2">魔族</Option>
            <Option value="0">其它</Option>
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

  getStatData(data) {
    let tian = 0;
    let mo = 0;
    let other = 0;
    data.forEach(v => {
      switch (v.type) {
        case 0:
          other++
          break
        case 1:
          tian++
          break
        case 2:
          mo++
          break
      }
    })
    return {tian, mo, other}
  }

  render() {
    const {playerList, loading} = this.props
    const statData = this.getStatData(playerList)
    return (
      <PageContainer>
        <Card extra={this.searchForm()}>
          <Row gutter={12}>
            <Col span={6}>
              <Card title="概况">
                <Row gutter={24}>
                  <Col span={12}>
                    <Statistic title="总数" value={playerList.length} style={{padding: "24px"}}
                               valueStyle={{color: "red"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="天族" value={statData.tian} style={{padding: "24px"}}
                               valueStyle={{color: "green"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="其它" value={statData.other} style={{padding: "24px"}}
                               valueStyle={{color: "orange"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="魔族" value={statData.mo} style={{padding: "24px"}} valueStyle={{color: "blue"}}/>
                  </Col>
                </Row>
              </Card>
            </Col>
            <Col span={18}>
              <Table
                bordered
                size="small"
                columns={this.columns}
                dataSource={playerList}
                rowKey={(record) => {
                  return record.name + record.type
                }}
                pagination={{
                  defaultPageSize: 15,
                  hideOnSinglePage: true,
                  showTotal: (total) => `共${total}条记录`,
                }}
                loading={loading}
              />
            </Col>
          </Row>
        </Card>
      </PageContainer>
    );
  }
}

export default Player;
