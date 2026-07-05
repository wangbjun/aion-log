import {Alert, Card, Col, notification, Progress, Row, Steps, Tag, Upload} from 'antd';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DatabaseOutlined,
  FileTextOutlined,
  InboxOutlined,
  LoadingOutlined,
} from '@ant-design/icons';
import {PageContainer} from '@ant-design/pro-layout';
import React from 'react';
import {importLogData, queryImportStatus} from '@/services/api';

const {Dragger} = Upload;
const {Step} = Steps;
class ImportLog extends React.Component {
  state = {
    importing: false,
    importStatus: 'idle',
    phase: 'idle',
    progress: 0,
    fileName: '',
    fileSize: '',
    fileLocation: '',
    statusMessage: '等待选择 Chat.log 文件',
    startedAt: '',
    finishedAt: '',
  }

  uploadRequestPending = false

  componentDidMount() {
    this.refreshImportStatus()
  }

  componentWillUnmount() {
    this.stopStatusPolling()
  }

  formatFileSize(size) {
    if (!size && size !== 0) {
      return ''
    }
    if (size < 1024) {
      return `${size} B`
    }
    if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(1)} KB`
    }
    if (size < 1024 * 1024 * 1024) {
      return `${(size / 1024 / 1024).toFixed(1)} MB`
    }
    return `${(size / 1024 / 1024 / 1024).toFixed(1)} GB`
  }

  getFileLocation(file) {
    return file.webkitRelativePath || '本地文件（浏览器不暴露完整路径）'
  }

  startStatusPolling = () => {
    this.stopStatusPolling()
    this.refreshImportStatus()
    this.statusTimer = setInterval(this.refreshImportStatus, 1000)
  }

  stopStatusPolling = () => {
    if (this.statusTimer) {
      clearInterval(this.statusTimer)
      this.statusTimer = null
    }
  }

  refreshImportStatus = async () => {
    try {
      const result = await queryImportStatus()
      if (result.code !== 200 || !result.data) {
        return
      }
      this.applyImportStatus(result.data)
    } catch (e) {
      // 状态轮询失败不打断导入请求，避免重复弹错误提示。
    }
  }

  applyImportStatus(status) {
    const statusValue = status.status || 'idle'
    if (this.uploadRequestPending && statusValue !== 'importing') {
      return
    }
    const importing = statusValue === 'importing'
    const nextState = {
      importing,
      importStatus: statusValue,
      phase: status.phase || 'idle',
      progress: status.progress || 0,
      fileName: status.fileName || this.state.fileName,
      fileSize: status.fileSize ? this.formatFileSize(status.fileSize) : this.state.fileSize,
      statusMessage: status.error || status.message || this.state.statusMessage,
      startedAt: status.startedAt || '',
      finishedAt: status.finishedAt || '',
    }
    if (status.fileName && !this.state.fileLocation) {
      nextState.fileLocation = '服务端已接收文件流'
    }
    this.setState(nextState)
    if (importing && !this.statusTimer) {
      this.statusTimer = setInterval(this.refreshImportStatus, 1000)
    } else if (!importing) {
      this.stopStatusPolling()
    }
  }


  importLog = async ({file, onSuccess, onError}) => {
    this.uploadRequestPending = true
    this.setState({
      importing: true,
      importStatus: 'importing',
      phase: 'upload',
      progress: 5,
      fileName: file.name,
      fileSize: this.formatFileSize(file.size),
      fileLocation: this.getFileLocation(file),
      statusMessage: '正在上传并解析日志，请不要关闭页面',
      startedAt: '',
      finishedAt: '',
    })
    this.startStatusPolling()
    try {
      const result = await importLogData(file)
      this.uploadRequestPending = false
      if (result.code !== 200) {
        await this.refreshImportStatus()
        notification.error({message: result.msg})
        onError && onError(new Error(result.msg))
        return
      }
      await this.refreshImportStatus()
      notification.success({message: '导入完成'})
      onSuccess && onSuccess(result, file)
    } catch (e) {
      this.uploadRequestPending = false
      this.setState({
          importing: false,
          importStatus: 'error',
          statusMessage: e.message || '导入失败',
        })
      this.stopStatusPolling()
      notification.error({message: e.message})
      onError && onError(e)
    }
  }

  getStatusMeta() {
    const {importStatus, phase, progress} = this.state
    const stepCurrent = ['upload', 'reset', 'parse'].includes(phase) ? 1 : 2
    if (importStatus === 'importing') {
      return {
        icon: <LoadingOutlined/>,
        title: '正在导入',
        tag: <Tag color="processing">导入中</Tag>,
        alertType: 'warning',
        progressStatus: 'active',
        percent: progress || 5,
        stepCurrent,
        stepStatus: 'process',
      }
    }
    if (importStatus === 'success') {
      return {
        icon: <CheckCircleOutlined/>,
        title: '导入完成',
        tag: <Tag color="success">已完成</Tag>,
        alertType: 'success',
        progressStatus: 'success',
        percent: 100,
        stepCurrent: 2,
        stepStatus: 'finish',
      }
    }
    if (importStatus === 'error') {
      return {
        icon: <CloseCircleOutlined/>,
        title: '导入失败',
        tag: <Tag color="error">失败</Tag>,
        alertType: 'error',
        progressStatus: 'exception',
        percent: 100,
        stepCurrent: 1,
        stepStatus: 'error',
      }
    }
    return {
      icon: <InboxOutlined/>,
      title: '等待文件',
      tag: <Tag>未开始</Tag>,
      alertType: 'info',
      progressStatus: 'normal',
      percent: 0,
      stepCurrent: 0,
      stepStatus: 'wait',
    }
  }

  renderStatusCard() {
    const {fileName, fileLocation, fileSize, finishedAt, startedAt, statusMessage} = this.state
    const statusMeta = this.getStatusMeta()
    const displayTime = finishedAt || startedAt || '-'

    return (
      <Card className={`battle-section-card battle-import-status battle-import-status-${this.state.importStatus}`}>
        <div className="battle-import-status-header">
          <div className="battle-import-status-title">
            <span className="battle-import-status-icon">{statusMeta.icon}</span>
            <span>{statusMeta.title}</span>
          </div>
          {statusMeta.tag}
        </div>
        <Progress
          percent={statusMeta.percent}
          status={statusMeta.progressStatus}
          showInfo={false}
        />
        <div className="battle-import-file">
          <span className="battle-muted">文件</span>
          <span>{fileName || '-'}</span>
        </div>
        <div className="battle-import-file">
          <span className="battle-muted">大小</span>
          <span>{fileSize || '-'}</span>
        </div>
        <div className="battle-import-file">
          <span className="battle-muted">位置</span>
          <span title={fileLocation || ''}>{fileLocation || '-'}</span>
        </div>
        <div className="battle-import-file">
          <span className="battle-muted">时间</span>
          <span>{displayTime}</span>
        </div>
        <Alert
          type={statusMeta.alertType}
          showIcon
          message={statusMessage}
        />
      </Card>
    )
  }

  render() {
    const {importing} = this.state
    const statusMeta = this.getStatusMeta()

    return (
      <PageContainer title={false}>
        <Row gutter={[16, 16]} className="battle-import-page">
          <Col xs={24} lg={15}>
            <Card className="battle-section-card battle-upload-panel battle-import-main">
              <div className="battle-import-main-header">
                <div>
                  <div className="battle-import-title">战斗日志导入</div>
                  <div className="battle-muted">Chat.log / TXT</div>
                </div>
                <Tag color={importing ? 'processing' : 'blue'}>{importing ? '处理中' : '可导入'}</Tag>
              </div>
              <Dragger
                accept=".log,.txt"
                customRequest={this.importLog}
                disabled={importing}
                maxCount={1}
                showUploadList={false}
              >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined/>
                </p>
                <p className="ant-upload-text">{importing ? '正在处理文件' : '拖拽 Chat.log 到这里'}</p>
                <p className="ant-upload-hint">
                  {importing ? '导入期间保持页面打开' : '点击选择文件也可以'}
                </p>
              </Dragger>
              <Row gutter={[12, 12]} className="battle-import-summary">
                <Col xs={24} md={8}>
                  <div className="battle-import-summary-item">
                    <FileTextOutlined/>
                    <span>读取日志</span>
                  </div>
                </Col>
                <Col xs={24} md={8}>
                  <div className="battle-import-summary-item">
                    <DatabaseOutlined/>
                    <span>重建数据</span>
                  </div>
                </Col>
                <Col xs={24} md={8}>
                  <div className="battle-import-summary-item">
                    <CheckCircleOutlined/>
                    <span>刷新统计</span>
                  </div>
                </Col>
              </Row>
              <Alert
                className="battle-import-note"
                type="info"
                showIcon
                message="Chatlog 只记录客户端附近的战斗行为"
                description="同一次战斗中，超出记录范围或网络延迟导致的缺失都可能影响统计结果。"
              />
            </Card>
          </Col>
          <Col xs={24} lg={9}>
            {this.renderStatusCard()}
            <Card title="处理流程" className="battle-section-card battle-import-flow">
              <Steps
                direction="vertical"
                size="small"
                current={statusMeta.stepCurrent}
                status={statusMeta.stepStatus}
              >
                <Step title="上传文件" description="支持 .log 和 .txt 格式"/>
                <Step title="解析战斗事件" description="提取玩家、技能、目标、伤害和击杀信息"/>
                <Step title="刷新统计结果" description="生成总览、日志、排行和时间线"/>
              </Steps>
            </Card>
          </Col>
        </Row>
      </PageContainer>
    );
  }
}

export default ImportLog;
