#include "ELRClient.h"
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>

ELRClient::ELRClient(QObject *parent) : QObject(parent)
{
    networkManager = new QNetworkAccessManager(this);
    monitoringTimer = new QTimer(this);
    elrApiBaseUrl = "http://localhost:8080/api";
    
    connect(monitoringTimer, &QTimer::timeout, this, &ELRClient::checkELRStatus);
    connect(networkManager, &QNetworkAccessManager::finished, this, [=](QNetworkReply *reply) {
        if (reply->url().toString().contains("status")) {
            onStatusReplyFinished();
        } else if (reply->url().toString().contains("containers")) {
            onContainersReplyFinished();
        }
    });
}

ELRClient::~ELRClient()
{
    delete networkManager;
    delete monitoringTimer;
}

void ELRClient::startMonitoring()
{
    monitoringTimer->start(5000); // 每5秒检查一次
    checkELRStatus(); // 立即检查一次
}

void ELRClient::stopMonitoring()
{
    monitoringTimer->stop();
}

void ELRClient::checkELRStatus()
{
    getELRStatus();
    getELRContainers();
}

void ELRClient::getELRStatus()
{
    QUrl url(elrApiBaseUrl + "/status");
    QNetworkRequest request(url);
    networkManager->get(request);
}

void ELRClient::getELRContainers()
{
    QUrl url(elrApiBaseUrl + "/containers");
    QNetworkRequest request(url);
    networkManager->get(request);
}

void ELRClient::onStatusReplyFinished()
{
    QNetworkReply *reply = qobject_cast<QNetworkReply*>(sender());
    if (!reply) return;
    
    if (reply->error() == QNetworkReply::NoError) {
        QByteArray responseData = reply->readAll();
        QJsonDocument jsonDoc = QJsonDocument::fromJson(responseData);
        QJsonObject jsonObj = jsonDoc.object();
        
        if (jsonObj.contains("status")) {
            QString status = jsonObj["status"].toString();
            emit statusUpdated(status);
        }
    } else {
        emit statusUpdated("未连接");
    }
    
    reply->deleteLater();
}

void ELRClient::onContainersReplyFinished()
{
    QNetworkReply *reply = qobject_cast<QNetworkReply*>(sender());
    if (!reply) return;
    
    QList<QMap<QString, QString>> containers;
    
    if (reply->error() == QNetworkReply::NoError) {
        QByteArray responseData = reply->readAll();
        QJsonDocument jsonDoc = QJsonDocument::fromJson(responseData);
        QJsonArray jsonArray = jsonDoc.array();
        
        for (int i = 0; i < jsonArray.size(); ++i) {
            QJsonObject containerObj = jsonArray[i].toObject();
            QMap<QString, QString> container;
            container["name"] = containerObj["name"].toString();
            container["status"] = containerObj["status"].toString();
            containers.append(container);
        }
    }
    
    emit containersUpdated(containers);
    reply->deleteLater();
}
